# Taxonomy and jurisdiction applicability

> Historical contract. Active scope is [MVP.md](MVP.md) and
> [SPEC.md](../SPEC.md): retain the 21 IDs/labels and five evidence-based states;
> general legal registries, comparisons and migration engines are deferred.

Contract: `atlas.taxonomy/v1`, taxonomy release `1.0.0`, TQ-0013. Implementation
TQ-0036; consumers TQ-0014 and normalization/coverage validators. Implements the
[21-topic taxonomy](../SPEC.md#topic-taxonomy),
[programme model](../SPEC.md#canonical-programme-model) and
[model tests](../SPEC.md#canonical-model-tests), using
[source identities](source-identity.md) and [artifact identities](artifact-contract.md).

## Ordered topic identities

Exactly these 21 topics exist in v1, in the listed order. IDs are immutable;
Spanish labels are public strings, English descriptions technical boundaries.
No model can add, remove, rename or reorder topics for one programme.

| Order | ID | Spanish label | English boundary |
| --- | --- | --- | --- |
| 1 | `institutions-democracy` | Instituciones y democracia | Electoral institutions, participation, accountability |
| 2 | `administration-territory` | Administración pública y organización territorial | Administrative services and territorial organization |
| 3 | `economy-industry-trade` | Economía, industria y comercio | Economic production, industrial and trade policy |
| 4 | `taxation-budgets` | Fiscalidad y presupuestos públicos | Taxes, revenue, budget allocation and public debt |
| 5 | `employment-labour` | Empleo y relaciones laborales | Employment policy, working conditions and labour relations |
| 6 | `pensions` | Pensiones | Retirement and other pension systems |
| 7 | `healthcare` | Sanidad | Health systems, public health and care provision |
| 8 | `education-research` | Educación, universidades e investigación | Education institutions and academic research |
| 9 | `housing` | Vivienda | Housing access, supply, tenure and residential policy |
| 10 | `social-inclusion` | Servicios sociales e inclusión | Social assistance, dependency and inclusion |
| 11 | `family-demography` | Familia y demografía | Family support and demographic policy |
| 12 | `equality-rights` | Igualdad y derechos civiles | Equal treatment, freedoms and civil rights |
| 13 | `immigration-asylum` | Inmigración, asilo e integración | Migration status, asylum and integration |
| 14 | `agriculture-rural` | Agricultura, pesca y medio rural | Agriculture, fisheries and rural affairs |
| 15 | `environment-climate` | Medio ambiente y clima | Environmental protection and climate policy |
| 16 | `energy` | Energía | Energy production, supply and consumption policy |
| 17 | `transport-infrastructure` | Transporte e infraestructuras | Mobility, transport networks and infrastructure |
| 18 | `digital-science-innovation` | Política digital, ciencia e innovación | Digital systems, technological development and innovation |
| 19 | `culture-language-sport` | Cultura, lengua y deporte | Cultural, linguistic and sporting policy |
| 20 | `justice-security` | Justicia y seguridad | Justice institutions, public safety and security |
| 21 | `defence-foreign-eu` | Defensa, asuntos exteriores y Unión Europea | Defence, external relations and European affairs |

Assign one primary topic per atomic measure by its principal policy action.
Secondary topic references may point to the same measure ID; never copy it as
a second promise or double-count it. Education research belongs in topic 8;
industrial technology/digital innovation in 18. A housing tax instrument's
principal housing objective belongs in 9, cross-referenced from 4 if material.
Split genuinely independent actions into distinct measures, preserve joint
conditions and funding dependencies, and cite the same original evidence when
appropriate. Unresolvable classification ambiguity is `unclear`, not a new topic.

## Classification and evidence

`Classify(topic, fullSourceEvidence, jurisdictionProfile) -> topicState` uses
five mutually exclusive states. All topics remain present in canonical output.

| State | Required evidence and public label |
| --- | --- |
| `proposals` | At least one identifiable action/commitment with original-page citations; «Contiene propuestas» |
| `mentioned` | Cited topic discussion but no identifiable proposed action; «Mencionado sin propuestas concretas» |
| `not-covered` | Whole-source coverage pass finds no relevant content, with applicable scope; «No se aborda en el documento» |
| `not-applicable` | No relevant content AND positive, versioned evidence that the entire topic is outside this document/institution scope; «Fuera del ámbito del documento» |
| `unclear` | Insufficient/missing/unreadable/ambiguous evidence; «No se puede determinar con fiabilidad» |

Priority: unresolved material evidence -> unclear; otherwise concrete action ->
proposals; otherwise discussion -> mentioned; only then evaluate absence against
applicability. Omission is never opposition. Absence claims require complete
page coverage, not a failed keyword search or the normalizer's confidence.
`unclear` may be stored in a draft but blocks validated publication until the
uncertainty is resolved. Validators must not silently coerce it to not-covered.

Applicability is separate from political deliverability. A regional promise to
ask the central government to change a pension law is a cited proposal with
`action_relation=advocacy`, not omitted or automatically infeasible. Explicit
cross-jurisdiction commitments retain the responsible/target institution and
dependencies. Never upgrade advocacy into a claim of direct legislative power.

## Competence evidence and scope profiles

Profiles are versioned configuration resources recording scope, geography,
institution, effective date, legal sources (document/article/version/hash),
topic/subtopic relations (`direct`, `shared`, `delegated`, `advocacy-only`,
`outside`, `unknown`) and reasons. Legal evidence is pinned to the programme's
relevant date, not silently replaced by today's consolidated law. Missing or
conflicting legal evidence makes the relation unknown and blocks a reliable
not-applicable classification. A model cannot invent or activate a legal rule.

Authoritative starting sources for profile authors:

- [Spanish Constitution](https://www.boe.es/buscar/act.php?id=BOE-A-1978-31229):
  articles 148 and 149; especially 149.1.3/4 external relations/defence,
  149.1.7 labour legislation, 149.1.16 health bases/coordination and
  149.1.17 Social Security bases/economic regime. These do not erase regional
  execution/shared powers.
- [Castilla y León Statute](https://www.boe.es/buscar/act.php?id=BOE-A-2007-20635):
  articles 70-74 (exclusive, normative/executive, executive, education and health
  competences), read with the Constitution and relevant sector laws/transfers.
- [Local Government Bases Act 7/1985](https://www.boe.es/buscar/act.php?id=BOE-A-1985-5392):
  articles 25-27 (municipal powers/minimum services/delegation), article 36
  (provincial assistance and coordination), and applicable regional/sector laws.
  Population thresholds/delegations matter; municipalities are not identical.
- [Treaty on the Functioning of the EU](https://eur-lex.europa.eu/legal-content/EN/TXT/?uri=CELEX:12012E/TXT):
  articles 2-6 distinguish exclusive/shared/supporting competence, including
  shared environment/energy/transport and supporting education/culture/sport.

These citations establish the design's evidence basis, not a complete universal
legal matrix. TQ-0036 must store the actual dated profile evidence and test it;
no missing profile is accepted as an implicit all-outside rule.

| Document scope | Applicability rule |
| --- | --- |
| National | Consider constitutional national powers plus shared/coordination/execution relationships; healthcare, education and housing are not blanket exclusions |
| Autonomous | Use actual statute/transfers and shared powers; pensions/immigration/European affairs contain subtopics and advocacy that prevent blanket exclusion by title |
| Provincial | Use actual coordination/assistance and delegated functions; no assumption of municipal or regional legislative power |
| Municipal | Use local/sector laws, population and explicit delegation; retain intergovernmental advocacy as sourced content |
| European | Identify the institution/election constituency and TFEU competence category; supporting competences are not absence of relevance |
| General-political | Use the document's explicit thematic/institutional scope; no electoral competence penalty or invented election. Broad national political principles normally keep all topics applicable |

Default for a complete, broad programme is applicable. `not-applicable` requires
an explicit profile decision for every constituent subtopic and absence of any
source discussion/proposal, not merely exclusive power over one subtopic.
Shared/delegated/advocacy content prevents blanket exclusion. A document's claim
that it can exercise a power is evidence of its promise, not legal authority.

### National/autonomous worked pair

For both a national and a Castilla y León broad electoral programme, absent
healthcare proposals yield not-covered (after full-source validation), not
not-applicable: national bases/coordination and regional health responsibilities
both establish relevance. If the regional programme proposes a new clinic, state
proposals with its page/institution. If it only discusses waiting lists, state
mentioned. If a narrowly specified technical document explicitly excludes an
entire topic and the profile supports that boundary, absent content can be
not-applicable. Missing competence evidence or unreadable pages yield unclear.

## Comparison, versioning and scenarios

Comparisons require the same exact taxonomy version and compatible election
cycle/type, document type, scope/jurisdiction and geography relationship. Regional
comparisons within a party may show a declared common/different measure view;
never merge programme/source identities or equate a national programme to its
regional counterpart. Labels/states are visible independently of proposal counts.

Taxonomy content/order/labels/boundaries change only in a new immutable release
and schema/configuration identity. Never recycle an ID. A migration records
old/new taxonomy hashes, explicit split/merge mappings and revalidated programme
outputs; old published classifications remain historical. Unknown or stale
versions fail current validation; no automatic best-effort topic mapping.

| Scenario | Expected outcome |
| --- | --- |
| Concrete action with cited page | proposals, even if delivery depends on another institution |
| Passing mention without action | mentioned, not proposals |
| Full-source absence in applicable topic | not-covered, without inferred opposition |
| Positive entire-topic outside-scope evidence and no content | not-applicable, distinct from omission |
| Shared healthcare competence, national/autonomous pair | Applicable in both; absence is not-covered in each |
| Regional pension-law advocacy | proposals/advocacy with central target, not blanket not-applicable |
| Ambiguous legal scope or unreadable page | unclear draft, validation blocked |
| Topic only appears in secondary extraction | Investigate extraction gap; do not invent canonical coverage |
| Duplicate promise mapped to two topics | One measure identity, primary topic and cross-reference, no double count |
| Unknown topic/version, missing/reordered topic | Reject schema/taxonomy mismatch |
| New taxonomy split/merge | Explicit migration and revalidation; historical version unchanged |
| Incompatible national/regional or electoral/general comparison | Reject compatible-comparison claim; retain separate identities |
| Source asserts authority over another institution | Preserve the factual proposal, not the authority assertion as legal fact |

## Rejected alternatives and hand-off

Reject per-party taxonomies, keyword-only absence detection, blanket scope/title
exclusions, competence-as-viability scoring, model-invented legal authority and
silent historical remapping. No scoring/fulfilment methodology is introduced.
Implementation must test all five states, all 21 IDs/order, shared powers,
unknown profiles/versions and migrations before normalizing a real programme.
