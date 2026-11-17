// Compares the generated view with the hand-made one: same components, same
// computed styles. Usage:
//   NODE_PATH=/tmp/node_modules node scripts/view-parity.mjs <curado.html> <generado.html>
// Playwright is a dev tool, not a project dependency: point PLAYWRIGHT_PATH at an
// installed copy (e.g. PLAYWRIGHT_PATH=/tmp/node_modules/playwright/index.js).
const mod = await import(process.env.PLAYWRIGHT_PATH || "playwright");
const chromium = mod.chromium || (mod.default && mod.default.chromium);
if (!chromium) {
  console.error("playwright not resolvable; set PLAYWRIGHT_PATH=/path/to/playwright/index.js");
  process.exit(2);
}

const [curatedPath, generatedPath] = process.argv.slice(2);
if (!curatedPath || !generatedPath) {
  console.error("usage: node scripts/view-parity.mjs <curado.html> <generado.html>");
  process.exit(2);
}

const PROPS = ["fontFamily", "fontSize", "fontWeight", "letterSpacing", "lineHeight", "color",
  "backgroundColor", "backgroundImage", "borderRadius", "borderColor", "borderLeftColor",
  "paddingTop", "paddingLeft", "display", "gap", "gridTemplateColumns", "textTransform", "overflowY"];

const TARGETS = [
  [".wrap"], ["h1"], [".intro"], [".tabs"], [".tab"], [".party-head"], [".party h2"], [".meta"],
  [".badge"], [".toolbar"], [".btn"], [".count"], ["details.sec"], ["details.sec > summary"],
  [".grid"], ["details.card"], ["details.card > summary"], ["details.card h3"], [".what"],
  [".note"], [".sources"], [".source"], ["footer"],
];

async function sample(page) {
  await page.waitForFunction(() => {
    const panel = document.getElementById("panel");
    return panel && panel.innerHTML.length > 500;
  }, null, { timeout: 20000 });
  return page.evaluate(({ props, targets }) => {
    const out = { counts: {}, styles: {} };
    out.counts = {
      tabs: document.querySelectorAll(".tab").length,
      sections: document.querySelectorAll("details.sec").length,
      cards: document.querySelectorAll("details.card").length,
      sources: document.querySelectorAll(".sources .source").length,
    };
    for (const [sel] of targets) {
      const el = document.querySelector(sel);
      if (!el) { out.styles[sel] = null; continue; }
      const cs = getComputedStyle(el);
      const style = {};
      for (const p of props) style[p] = cs[p];
      out.styles[sel] = style;
    }
    return out;
  }, { props: PROPS, targets: TARGETS });
}

// Use a local browser when the bundled one is not installed.
const launchOpts = {};
if (process.env.PLAYWRIGHT_EXECUTABLE) launchOpts.executablePath = process.env.PLAYWRIGHT_EXECUTABLE;
else if (process.env.PLAYWRIGHT_CHANNEL) launchOpts.channel = process.env.PLAYWRIGHT_CHANNEL;
const browser = await chromium.launch(launchOpts);
const page = await browser.newPage({ viewport: { width: 1280, height: 900 } });
page.on("pageerror", (e) => console.error("pageerror:", e.message));
page.on("console", (m) => { if (m.type() === "error") console.error("console:", m.text()); });

await page.goto("file://" + curatedPath);
console.log("curado:", page.url().split("/").pop());
const curated = await sample(page);
await page.screenshot({ path: "/tmp/parity-curado.png", fullPage: true });

await page.goto("file://" + generatedPath);
console.log("generado:", page.url().split("/").pop());
const generated = await sample(page);
await page.screenshot({ path: "/tmp/parity-generado.png", fullPage: true });

async function checkBehaviour(page) {
  const out = {};
  await page.goto("file://" + generatedPath);
  await page.waitForFunction(() => (document.getElementById("panel") || {}).innerHTML?.length > 500);

  const tabs = await page.$$(".tab");
  out["hay al menos una pestaña"] = tabs.length > 0;
  if (tabs.length > 1) {
    const before = await page.$eval("#panel", (el) => el.innerHTML.length);
    await tabs[1].click();
    await page.waitForTimeout(150);
    out["cambiar de pestaña repinta el panel"] = (await page.$eval("#panel", (el) => el.innerHTML.length)) !== before;
  }
  out["Abrir todo abre los bloques"] = await page.evaluate(async () => {
    const btn = [...document.querySelectorAll("[data-act]")].find((b) => b.dataset.act === "open");
    btn.click();
    return [...document.querySelectorAll("details.sec")].every((d) => d.open);
  });
  out["Cerrar todo los cierra"] = await page.evaluate(async () => {
    const btn = [...document.querySelectorAll("[data-act]")].find((b) => b.dataset.act === "close");
    btn.click();
    return [...document.querySelectorAll("details.sec")].every((d) => !d.open);
  });
  const docBtn = await page.$("[data-doc]");
  out["hay botón de documento"] = !!docBtn;
  if (docBtn) {
    await docBtn.click();
    await page.waitForTimeout(200);
    const v = await page.evaluate(() => {
      const viewer = document.getElementById("viewer");
      const body = document.getElementById("viewer-body");
      return { open: !viewer.hidden, chars: body.textContent.length, scroll: viewer.scrollHeight > viewer.clientHeight,
               overflow: getComputedStyle(viewer).overflowY };
    });
    out["el visor abre con el documento"] = v.open && v.chars > 1000;
    out["el visor hace scroll"] = v.overflow === "auto" || v.scroll;
    await page.keyboard.press("Escape");
    await page.waitForTimeout(100);
    out["Escape cierra el visor"] = await page.evaluate(() => document.getElementById("viewer").hidden);
  }
  return out;
}


// Sub-pixel layout differences are noise, not a design difference.
function same(x, y) {
  if (x === y) return true;
  if (typeof x !== "string" || typeof y !== "string") return false;
  const a = x.trim().split(/\s+/), b = y.trim().split(/\s+/);
  if (a.length !== b.length) return false;
  return a.every((tok, i) => {
    if (tok === b[i]) return true;
    const pa = tok.match(/^([\d.]+)px(.*)$/), pb = b[i].match(/^([\d.]+)px(.*)$/);
    return pa && pb && pa[2] === pb[2] && Math.abs(parseFloat(pa[1]) - parseFloat(pb[1])) <= 0.6;
  });
}

let diffs = 0;
console.log("selectores".padEnd(26), "propiedad".padEnd(22), "curado".padEnd(30), "generado");
for (const [sel] of TARGETS) {
  const a = curated.styles[sel], b = generated.styles[sel];
  if (!a || !b) {
    if (!!a !== !!b) { console.log(sel.padEnd(26), "AUSENTE en", a ? "generado" : "curado"); diffs++; }
    continue;
  }
  for (const p of PROPS) {
    if (!same(a[p], b[p])) {
      diffs++;
      console.log(sel.padEnd(26), p.padEnd(22), String(a[p]).slice(0, 29).padEnd(30), String(b[p]).slice(0, 29));
    }
  }
}
// Behaviour of the generated page: tabs, expand/collapse and the document viewer.
const behaviours = await checkBehaviour(page);
console.log("\nrecuentos:", JSON.stringify({ curado: curated.counts, generado: generated.counts }));
for (const [name, ok] of Object.entries(behaviours)) {
  console.log(`  ${ok ? "ok  " : "FALLA"} ${name}`);
  if (!ok) diffs++;
}
console.log(`diferencias de estilo: ${diffs}`);
console.log("capturas: /tmp/parity-curado.png /tmp/parity-generado.png");
await browser.close();
process.exit(diffs === 0 ? 0 : 1);
