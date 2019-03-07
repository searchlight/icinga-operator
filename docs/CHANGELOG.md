---
title: Changelog | Searchlight
description: Changelog
menu:
  product_searchlight_8.0.0-rc.0:
    identifier: changelog-searchlight
    name: Changelog
    parent: welcome
    weight: 10
product_name: searchlight
menu_name: product_searchlight_8.0.0-rc.0
section_menu_id: welcome
url: /products/searchlight/8.0.0-rc.0/welcome/changelog/
aliases:
  - /products/searchlight/8.0.0-rc.0/CHANGELOG/
---

# Change Log

## [8.0.0-rc.0](https://github.com/appscode/searchlight/tree/8.0.0-rc.0) (2019-03-07)
[Full Changelog](https://github.com/appscode/searchlight/compare/7.0.0...8.0.0-rc.0)

**Fixed bugs:**

- Unknown : var "icinga.checkInterval" is not registered in SearchlightPlugin [\#397](https://github.com/appscode/searchlight/issues/397)
- Correctly handle ignored openapi prefixes [\#398](https://github.com/appscode/searchlight/pull/398) ([tamalsaha](https://github.com/tamalsaha))

**Closed issues:**

- Helm chart is not operable [\#419](https://github.com/appscode/searchlight/issues/419)

**Merged pull requests:**

- Use onessl 0.10.0 [\#437](https://github.com/appscode/searchlight/pull/437) ([tamalsaha](https://github.com/tamalsaha))
- Update Kubernetes client libraries to 1.13.0 [\#436](https://github.com/appscode/searchlight/pull/436) ([tamalsaha](https://github.com/tamalsaha))
- Various fixes to installer [\#435](https://github.com/appscode/searchlight/pull/435) ([tamalsaha](https://github.com/tamalsaha))
- Pin go-ini dependency [\#434](https://github.com/appscode/searchlight/pull/434) ([tamalsaha](https://github.com/tamalsaha))
- Update webhook error message format for Kubernetes 1.13+ [\#433](https://github.com/appscode/searchlight/pull/433) ([tamalsaha](https://github.com/tamalsaha))
- Remove openshift dependencies [\#431](https://github.com/appscode/searchlight/pull/431) ([tamalsaha](https://github.com/tamalsaha))
- Set periodic analytics [\#430](https://github.com/appscode/searchlight/pull/430) ([tamalsaha](https://github.com/tamalsaha))
- Update Kubernetes client libraries to 1.12.0 [\#429](https://github.com/appscode/searchlight/pull/429) ([tamalsaha](https://github.com/tamalsaha))
- Pass image pull secrets for cleaner job in chart [\#425](https://github.com/appscode/searchlight/pull/425) ([tamalsaha](https://github.com/tamalsaha))
- Update kubernetes client libraries to 1.12.0 [\#424](https://github.com/appscode/searchlight/pull/424) ([tamalsaha](https://github.com/tamalsaha))
- Document LogLevel for chart [\#423](https://github.com/appscode/searchlight/pull/423) ([tamalsaha](https://github.com/tamalsaha))
- Check if Kubernetes version is supported before running operator [\#422](https://github.com/appscode/searchlight/pull/422) ([tamalsaha](https://github.com/tamalsaha))
- Enable webhooks by default in chart [\#421](https://github.com/appscode/searchlight/pull/421) ([tamalsaha](https://github.com/tamalsaha))
- Fix service ports in chart [\#420](https://github.com/appscode/searchlight/pull/420) ([tamalsaha](https://github.com/tamalsaha))
- Fixed missing "end" in helm charts secret.yaml [\#418](https://github.com/appscode/searchlight/pull/418) ([maximilize](https://github.com/maximilize))
- Use kubernetes-1.11.3 [\#417](https://github.com/appscode/searchlight/pull/417) ([tamalsaha](https://github.com/tamalsaha))
- Update CertStore [\#416](https://github.com/appscode/searchlight/pull/416) ([tamalsaha](https://github.com/tamalsaha))
- Support pod annotations in chart [\#415](https://github.com/appscode/searchlight/pull/415) ([tamalsaha](https://github.com/tamalsaha))
- Set serviceAccount for clearner job [\#414](https://github.com/appscode/searchlight/pull/414) ([tamalsaha](https://github.com/tamalsaha))
- Cleanup webhooks when chart is deleted [\#413](https://github.com/appscode/searchlight/pull/413) ([tamalsaha](https://github.com/tamalsaha))
- Revendor [\#412](https://github.com/appscode/searchlight/pull/412) ([tamalsaha](https://github.com/tamalsaha))
- Update pipeline [\#411](https://github.com/appscode/searchlight/pull/411) ([tahsinrahman](https://github.com/tahsinrahman))
- Revendor api [\#410](https://github.com/appscode/searchlight/pull/410) ([tamalsaha](https://github.com/tamalsaha))
- fix uninstall for concourse [\#409](https://github.com/appscode/searchlight/pull/409) ([tahsinrahman](https://github.com/tahsinrahman))
- Fix chart values file [\#408](https://github.com/appscode/searchlight/pull/408) ([tamalsaha](https://github.com/tamalsaha))
- Improve Helm chart options [\#407](https://github.com/appscode/searchlight/pull/407) ([tamalsaha](https://github.com/tamalsaha))
- Use concourse scripts from libbuild [\#405](https://github.com/appscode/searchlight/pull/405) ([tahsinrahman](https://github.com/tahsinrahman))
- Add categories [\#404](https://github.com/appscode/searchlight/pull/404) ([tamalsaha](https://github.com/tamalsaha))
- Enable status sub resource for crd yamls [\#403](https://github.com/appscode/searchlight/pull/403) ([tamalsaha](https://github.com/tamalsaha))
- Move crds to api folder [\#402](https://github.com/appscode/searchlight/pull/402) ([tamalsaha](https://github.com/tamalsaha))
- Fix extended apiserver issues with Kubernetes 1.11 [\#401](https://github.com/appscode/searchlight/pull/401) ([tamalsaha](https://github.com/tamalsaha))
- fix CheckCommand generation for webhooks [\#399](https://github.com/appscode/searchlight/pull/399) ([mmta](https://github.com/mmta))
- Use version and additional columns for crds [\#396](https://github.com/appscode/searchlight/pull/396) ([tamalsaha](https://github.com/tamalsaha))
- Don't add admission/v1beta1 group as a prioritized version [\#395](https://github.com/appscode/searchlight/pull/395) ([tamalsaha](https://github.com/tamalsaha))
- Enable status subresource for voyager crds [\#394](https://github.com/appscode/searchlight/pull/394) ([tamalsaha](https://github.com/tamalsaha))
- Update client-go to v8.0.0 [\#393](https://github.com/appscode/searchlight/pull/393) ([tamalsaha](https://github.com/tamalsaha))
- Format shell scripts [\#392](https://github.com/appscode/searchlight/pull/392) ([tamalsaha](https://github.com/tamalsaha))
- Remove --ops-address flag [\#390](https://github.com/appscode/searchlight/pull/390) ([tamalsaha](https://github.com/tamalsaha))
- Update searchlight version in readme [\#389](https://github.com/appscode/searchlight/pull/389) ([tamalsaha](https://github.com/tamalsaha))
- Fix sample ack.yaml [\#388](https://github.com/appscode/searchlight/pull/388) ([tamalsaha](https://github.com/tamalsaha))
- Fix openapi schema for metav1.Duration [\#387](https://github.com/appscode/searchlight/pull/387) ([tamalsaha](https://github.com/tamalsaha))
- Move openapi-spec to api folder [\#386](https://github.com/appscode/searchlight/pull/386) ([tamalsaha](https://github.com/tamalsaha))
- Add togglable tabs for Installation: Script & Helm [\#385](https://github.com/appscode/searchlight/pull/385) ([sajibcse68](https://github.com/sajibcse68))
- Handle node status unknown for node-status checks [\#384](https://github.com/appscode/searchlight/pull/384) ([tamalsaha](https://github.com/tamalsaha))

## [7.0.0](https://github.com/appscode/searchlight/tree/7.0.0) (2018-06-02)
[Full Changelog](https://github.com/appscode/searchlight/compare/7.0.0-rc.0...7.0.0)

**Merged pull requests:**

- Prepare docs for 7.0.0 release [\#383](https://github.com/appscode/searchlight/pull/383) ([tamalsaha](https://github.com/tamalsaha))
- Revendor [\#382](https://github.com/appscode/searchlight/pull/382) ([tamalsaha](https://github.com/tamalsaha))
- Improve installer [\#381](https://github.com/appscode/searchlight/pull/381) ([tamalsaha](https://github.com/tamalsaha))
- concourse [\#379](https://github.com/appscode/searchlight/pull/379) ([tahsinrahman](https://github.com/tahsinrahman))
- Update changelog [\#378](https://github.com/appscode/searchlight/pull/378) ([tamalsaha](https://github.com/tamalsaha))

## [7.0.0-rc.0](https://github.com/appscode/searchlight/tree/7.0.0-rc.0) (2018-05-25)
[Full Changelog](https://github.com/appscode/searchlight/compare/5.1.1...7.0.0-rc.0)

**Implemented enhancements:**

- Check expiration for any cert [\#275](https://github.com/appscode/searchlight/issues/275)
- Support webhook based custom plugin [\#336](https://github.com/appscode/searchlight/pull/336) ([aerokite](https://github.com/aerokite))
- Add tests for plugins [\#313](https://github.com/appscode/searchlight/pull/313) ([aerokite](https://github.com/aerokite))

**Fixed bugs:**

- Permission missing for /healthz request in Rbac roles [\#351](https://github.com/appscode/searchlight/issues/351)
- Fix HTTP client with incluster config [\#166](https://github.com/appscode/searchlight/issues/166)

**Closed issues:**

- support removing acknowledgement [\#299](https://github.com/appscode/searchlight/issues/299)
- Sending custom notification from IcingaWeb2 does not work [\#297](https://github.com/appscode/searchlight/issues/297)
- Allow `pausing` alerts [\#295](https://github.com/appscode/searchlight/issues/295)
- Support user provided plugin [\#293](https://github.com/appscode/searchlight/issues/293)
- Record incidents and notifications via CRD/EAS [\#292](https://github.com/appscode/searchlight/issues/292)
- Add test for plugins [\#289](https://github.com/appscode/searchlight/issues/289)
- Replace cfssl with client-go cert utils [\#231](https://github.com/appscode/searchlight/issues/231)
- Add e2e test for notification [\#205](https://github.com/appscode/searchlight/issues/205)
- Handle Node conditions [\#198](https://github.com/appscode/searchlight/issues/198)
- Make http endpoint a UAS [\#98](https://github.com/appscode/searchlight/issues/98)

**Merged pull requests:**

- Refactor acknowledgement storage implementation [\#377](https://github.com/appscode/searchlight/pull/377) ([tamalsaha](https://github.com/tamalsaha))
- Use internal type to implement storage [\#376](https://github.com/appscode/searchlight/pull/376) ([tamalsaha](https://github.com/tamalsaha))
- Add api password for e2e-test [\#375](https://github.com/appscode/searchlight/pull/375) ([aerokite](https://github.com/aerokite))
- fix NotificationCommand arguments [\#374](https://github.com/appscode/searchlight/pull/374) ([aerokite](https://github.com/aerokite))
- provide pods/exec resource permission [\#373](https://github.com/appscode/searchlight/pull/373) ([aerokite](https://github.com/aerokite))
- add documentation for usage of stride notifier [\#372](https://github.com/appscode/searchlight/pull/372) ([aerokite](https://github.com/aerokite))
- create all built-in SearchlightPlugin at runtime [\#371](https://github.com/appscode/searchlight/pull/371) ([aerokite](https://github.com/aerokite))
- Fix chart [\#369](https://github.com/appscode/searchlight/pull/369) ([tamalsaha](https://github.com/tamalsaha))
- Various installer improvements [\#368](https://github.com/appscode/searchlight/pull/368) ([tamalsaha](https://github.com/tamalsaha))
- update webhook plugin [\#366](https://github.com/appscode/searchlight/pull/366) ([aerokite](https://github.com/aerokite))
- Revendor go-notify [\#365](https://github.com/appscode/searchlight/pull/365) ([tamalsaha](https://github.com/tamalsaha))
- sort arguments in CheckCommand [\#364](https://github.com/appscode/searchlight/pull/364) ([aerokite](https://github.com/aerokite))
- Don't panic if admission options is nil [\#363](https://github.com/appscode/searchlight/pull/363) ([tamalsaha](https://github.com/tamalsaha))
- Disable admission controllers for webhook server [\#362](https://github.com/appscode/searchlight/pull/362) ([tamalsaha](https://github.com/tamalsaha))
- Add Update\*\*\*Status helpers [\#361](https://github.com/appscode/searchlight/pull/361) ([tamalsaha](https://github.com/tamalsaha))
- Update to client-go 7.0.0 [\#360](https://github.com/appscode/searchlight/pull/360) ([tamalsaha](https://github.com/tamalsaha))
- Improve installer [\#359](https://github.com/appscode/searchlight/pull/359) ([tamalsaha](https://github.com/tamalsaha))
- add command for webhook plugin [\#357](https://github.com/appscode/searchlight/pull/357) ([aerokite](https://github.com/aerokite))
- Add documentation for custom plugin [\#356](https://github.com/appscode/searchlight/pull/356) ([aerokite](https://github.com/aerokite))
- Generate non-namespaced client for plugins [\#355](https://github.com/appscode/searchlight/pull/355) ([aerokite](https://github.com/aerokite))
- add patch permission [\#354](https://github.com/appscode/searchlight/pull/354) ([aerokite](https://github.com/aerokite))
- Fix docs [\#353](https://github.com/appscode/searchlight/pull/353) ([aerokite](https://github.com/aerokite))
- Various fixes to searchlight installer [\#352](https://github.com/appscode/searchlight/pull/352) ([tamalsaha](https://github.com/tamalsaha))
- Remove jessie icinga [\#349](https://github.com/appscode/searchlight/pull/349) ([aerokite](https://github.com/aerokite))
- Introduce properties for plugins vars [\#348](https://github.com/appscode/searchlight/pull/348) ([aerokite](https://github.com/aerokite))
- Migrate builtin check commands to Plugin crd [\#347](https://github.com/appscode/searchlight/pull/347) ([aerokite](https://github.com/aerokite))
- Add RBAC instructions for GKE cluster [\#346](https://github.com/appscode/searchlight/pull/346) ([tamalsaha](https://github.com/tamalsaha))
- Update chart repository location [\#345](https://github.com/appscode/searchlight/pull/345) ([tamalsaha](https://github.com/tamalsaha))
- Support installing from local installer scripts [\#344](https://github.com/appscode/searchlight/pull/344) ([tamalsaha](https://github.com/tamalsaha))
- Move swagger.json to openapi-spec folder [\#343](https://github.com/appscode/searchlight/pull/343) ([tamalsaha](https://github.com/tamalsaha))
- Regenerate swagger.json [\#342](https://github.com/appscode/searchlight/pull/342) ([tamalsaha](https://github.com/tamalsaha))
- Generate swagger.json [\#341](https://github.com/appscode/searchlight/pull/341) ([tamalsaha](https://github.com/tamalsaha))
- Add install pkg for crds [\#340](https://github.com/appscode/searchlight/pull/340) ([tamalsaha](https://github.com/tamalsaha))
- Skip setting ListKind [\#339](https://github.com/appscode/searchlight/pull/339) ([tamalsaha](https://github.com/tamalsaha))
- Add CRD Validation [\#338](https://github.com/appscode/searchlight/pull/338) ([tamalsaha](https://github.com/tamalsaha))
- Generate openapi spec [\#337](https://github.com/appscode/searchlight/pull/337) ([tamalsaha](https://github.com/tamalsaha))
- Fix install script for minikube 0.24.x \(Kube 1.8.0\) [\#335](https://github.com/appscode/searchlight/pull/335) ([tamalsaha](https://github.com/tamalsaha))
- Fix comment for LastNotificationType in IncidentStatus [\#334](https://github.com/appscode/searchlight/pull/334) ([aerokite](https://github.com/aerokite))
- fix typo [\#333](https://github.com/appscode/searchlight/pull/333) ([aerokite](https://github.com/aerokite))
- Garbage collect incidents older than 90 days [\#332](https://github.com/appscode/searchlight/pull/332) ([tamalsaha](https://github.com/tamalsaha))
- Document user roles [\#331](https://github.com/appscode/searchlight/pull/331) ([tamalsaha](https://github.com/tamalsaha))
- Update docs for json\_path [\#330](https://github.com/appscode/searchlight/pull/330) ([tamalsaha](https://github.com/tamalsaha))
- Correctly install validation webhook [\#329](https://github.com/appscode/searchlight/pull/329) ([tamalsaha](https://github.com/tamalsaha))
- Fix : No such file or directory: '$GOPATH/src/github.com/appscode/sea… [\#327](https://github.com/appscode/searchlight/pull/327) ([YangYongZhi](https://github.com/YangYongZhi))
- Add docs for adding check command [\#326](https://github.com/appscode/searchlight/pull/326) ([aerokite](https://github.com/aerokite))
- Fix build on mac [\#325](https://github.com/appscode/searchlight/pull/325) ([tamalsaha](https://github.com/tamalsaha))
- Skip downloading onessl is already exists [\#324](https://github.com/appscode/searchlight/pull/324) ([tamalsaha](https://github.com/tamalsaha))
- Fix installer script [\#323](https://github.com/appscode/searchlight/pull/323) ([tamalsaha](https://github.com/tamalsaha))
- Use server cert for icinga [\#322](https://github.com/appscode/searchlight/pull/322) ([tamalsaha](https://github.com/tamalsaha))
- Write auto-generated icinga certs to disk [\#321](https://github.com/appscode/searchlight/pull/321) ([tamalsaha](https://github.com/tamalsaha))
- Prepare docs for 7.0.0-rc.0 [\#320](https://github.com/appscode/searchlight/pull/320) ([tamalsaha](https://github.com/tamalsaha))
- Revendor kutil and jsonpatch library [\#319](https://github.com/appscode/searchlight/pull/319) ([tamalsaha](https://github.com/tamalsaha))
- Add changelog [\#318](https://github.com/appscode/searchlight/pull/318) ([tamalsaha](https://github.com/tamalsaha))
- Reorg objects deleted in uninstall command [\#317](https://github.com/appscode/searchlight/pull/317) ([tamalsaha](https://github.com/tamalsaha))
- Add tests for plugins [\#316](https://github.com/appscode/searchlight/pull/316) ([tamalsaha](https://github.com/tamalsaha))
- Add e2e test for notifier [\#315](https://github.com/appscode/searchlight/pull/315) ([aerokite](https://github.com/aerokite))
- Rename --analytics to --enable-analytics [\#314](https://github.com/appscode/searchlight/pull/314) ([tamalsaha](https://github.com/tamalsaha))
- send verbosity as Arg and analytics as Env [\#312](https://github.com/appscode/searchlight/pull/312) ([aerokite](https://github.com/aerokite))
- Revendor webhook api [\#311](https://github.com/appscode/searchlight/pull/311) ([tamalsaha](https://github.com/tamalsaha))
- update check\_json\_path [\#310](https://github.com/appscode/searchlight/pull/310) ([aerokite](https://github.com/aerokite))
- Update check\_node\_status to support other NodeCondition [\#309](https://github.com/appscode/searchlight/pull/309) ([aerokite](https://github.com/aerokite))
- Add check\_any\_cert plugin [\#307](https://github.com/appscode/searchlight/pull/307) ([aerokite](https://github.com/aerokite))
- Add incidents and Acknowledgements to user roles [\#306](https://github.com/appscode/searchlight/pull/306) ([tamalsaha](https://github.com/tamalsaha))
- Replace cfssl with client-go cert utils [\#305](https://github.com/appscode/searchlight/pull/305) ([tamalsaha](https://github.com/tamalsaha))
- Remove internal types for CRDs [\#304](https://github.com/appscode/searchlight/pull/304) ([tamalsaha](https://github.com/tamalsaha))
- support pausing check [\#303](https://github.com/appscode/searchlight/pull/303) ([aerokite](https://github.com/aerokite))
- Support delete acknowledgement [\#302](https://github.com/appscode/searchlight/pull/302) ([aerokite](https://github.com/aerokite))
- Fix build [\#301](https://github.com/appscode/searchlight/pull/301) ([tamalsaha](https://github.com/tamalsaha))
- Add travis.yml [\#300](https://github.com/appscode/searchlight/pull/300) ([tamalsaha](https://github.com/tamalsaha))
- Rename states to title case [\#298](https://github.com/appscode/searchlight/pull/298) ([tamalsaha](https://github.com/tamalsaha))
- Record incidents and notifications via CRD/EAS [\#296](https://github.com/appscode/searchlight/pull/296) ([tamalsaha](https://github.com/tamalsaha))
- Merge admission webhook and operator into one binary [\#291](https://github.com/appscode/searchlight/pull/291) ([tamalsaha](https://github.com/tamalsaha))
- Remove individual binaries for plugins [\#290](https://github.com/appscode/searchlight/pull/290) ([tamalsaha](https://github.com/tamalsaha))
- Update readme to 5.1.1 release [\#287](https://github.com/appscode/searchlight/pull/287) ([tamalsaha](https://github.com/tamalsaha))
- Use workqueue [\#230](https://github.com/appscode/searchlight/pull/230) ([tamalsaha](https://github.com/tamalsaha))

## [5.1.1](https://github.com/appscode/searchlight/tree/5.1.1) (2018-03-06)
[Full Changelog](https://github.com/appscode/searchlight/compare/5.1.0...5.1.1)

**Fixed bugs:**

- If no service available, delete host [\#285](https://github.com/appscode/searchlight/pull/285) ([aerokite](https://github.com/aerokite))

**Closed issues:**

- Migrate from 1.5.9 to 3.0.0 [\#194](https://github.com/appscode/searchlight/issues/194)
- Rewrite searchlight design doc [\#172](https://github.com/appscode/searchlight/issues/172)
- New alerts [\#169](https://github.com/appscode/searchlight/issues/169)
- Install as critical addon [\#118](https://github.com/appscode/searchlight/issues/118)
- Use webhook notifier for appscode api [\#94](https://github.com/appscode/searchlight/issues/94)

**Merged pull requests:**

- Prepare docs for 5.1.1 release [\#286](https://github.com/appscode/searchlight/pull/286) ([tamalsaha](https://github.com/tamalsaha))
- Make it clear that installer is a single command [\#284](https://github.com/appscode/searchlight/pull/284) ([tamalsaha](https://github.com/tamalsaha))
- Fix installer [\#283](https://github.com/appscode/searchlight/pull/283) ([tamalsaha](https://github.com/tamalsaha))
- Update chart to match RBAC best practices for charts [\#281](https://github.com/appscode/searchlight/pull/281) ([tamalsaha](https://github.com/tamalsaha))
- Document webhook notifier [\#280](https://github.com/appscode/searchlight/pull/280) ([tamalsaha](https://github.com/tamalsaha))
- Support --enable-admission-webhook=false [\#278](https://github.com/appscode/searchlight/pull/278) ([tamalsaha](https://github.com/tamalsaha))
- Support multiple webhooks of same apiversion [\#277](https://github.com/appscode/searchlight/pull/277) ([tamalsaha](https://github.com/tamalsaha))
- Sync chart to stable charts repo [\#276](https://github.com/appscode/searchlight/pull/276) ([tamalsaha](https://github.com/tamalsaha))
- Fix RBAC permission [\#274](https://github.com/appscode/searchlight/pull/274) ([tamalsaha](https://github.com/tamalsaha))
- Fix RBAC permission [\#273](https://github.com/appscode/searchlight/pull/273) ([tamalsaha](https://github.com/tamalsaha))
- Delete internal types [\#272](https://github.com/appscode/searchlight/pull/272) ([tamalsaha](https://github.com/tamalsaha))
- Use rbc/v1 apis [\#271](https://github.com/appscode/searchlight/pull/271) ([tamalsaha](https://github.com/tamalsaha))
- Create user facing aggregate roles [\#270](https://github.com/appscode/searchlight/pull/270) ([tamalsaha](https://github.com/tamalsaha))
- Use ${} form for onessl envsubst [\#269](https://github.com/appscode/searchlight/pull/269) ([tamalsaha](https://github.com/tamalsaha))
- Merge uninstall script into installer [\#268](https://github.com/appscode/searchlight/pull/268) ([tamalsaha](https://github.com/tamalsaha))
- Copy generic-admission-server into pkg [\#267](https://github.com/appscode/searchlight/pull/267) ([tamalsaha](https://github.com/tamalsaha))
- Cut 6.0.0-alpha.0 [\#265](https://github.com/appscode/searchlight/pull/265) ([tamalsaha](https://github.com/tamalsaha))
- Add ValidatingAdmissionWebhook for CRDs [\#264](https://github.com/appscode/searchlight/pull/264) ([tamalsaha](https://github.com/tamalsaha))
- Fix instructions for using private docker registry [\#263](https://github.com/appscode/searchlight/pull/263) ([tamalsaha](https://github.com/tamalsaha))
- Use installer script [\#262](https://github.com/appscode/searchlight/pull/262) ([tamalsaha](https://github.com/tamalsaha))
- Update client-go to v0.6.0 [\#261](https://github.com/appscode/searchlight/pull/261) ([tamalsaha](https://github.com/tamalsaha))
- Regenerate clients [\#259](https://github.com/appscode/searchlight/pull/259) ([tamalsaha](https://github.com/tamalsaha))

## [5.1.0](https://github.com/appscode/searchlight/tree/5.1.0) (2018-01-17)
[Full Changelog](https://github.com/appscode/searchlight/compare/5.0.0...5.1.0)

**Merged pull requests:**

- Prepare docs for 5.1.0 [\#258](https://github.com/appscode/searchlight/pull/258) ([tamalsaha](https://github.com/tamalsaha))
- Fix docs to make vars string=\>string [\#257](https://github.com/appscode/searchlight/pull/257) ([tamalsaha](https://github.com/tamalsaha))
- Support Telegram as notifier [\#256](https://github.com/appscode/searchlight/pull/256) ([tamalsaha](https://github.com/tamalsaha))

## [5.0.0](https://github.com/appscode/searchlight/tree/5.0.0) (2018-01-03)
[Full Changelog](https://github.com/appscode/searchlight/compare/4.0.0...5.0.0)

**Implemented enhancements:**

- Support hipchat server [\#237](https://github.com/appscode/searchlight/issues/237)

**Fixed bugs:**

- Failed to create events on CRD objects [\#216](https://github.com/appscode/searchlight/issues/216)

**Merged pull requests:**

- Fix analytics client id in GKE [\#255](https://github.com/appscode/searchlight/pull/255) ([tamalsaha](https://github.com/tamalsaha))
- Prepare docs for 5.0.0 [\#254](https://github.com/appscode/searchlight/pull/254) ([tamalsaha](https://github.com/tamalsaha))
- Revendor kutil [\#253](https://github.com/appscode/searchlight/pull/253) ([tamalsaha](https://github.com/tamalsaha))
- Reorganize docs [\#252](https://github.com/appscode/searchlight/pull/252) ([sajibcse68](https://github.com/sajibcse68))
- Support hipchat server [\#251](https://github.com/appscode/searchlight/pull/251) ([tamalsaha](https://github.com/tamalsaha))
- Remove TryPatch method [\#250](https://github.com/appscode/searchlight/pull/250) ([tamalsaha](https://github.com/tamalsaha))
- Indicate mutation in PATCH helper method return [\#249](https://github.com/appscode/searchlight/pull/249) ([tamalsaha](https://github.com/tamalsaha))
- Set analytics ClientID [\#247](https://github.com/appscode/searchlight/pull/247) ([tamalsaha](https://github.com/tamalsaha))
- Update gendocs script to generate front matter [\#245](https://github.com/appscode/searchlight/pull/245) ([tamalsaha](https://github.com/tamalsaha))
- Add front matter for reference/ [\#244](https://github.com/appscode/searchlight/pull/244) ([sajibcse68](https://github.com/sajibcse68))
- Fix section\_menu\_id for tutorials root files [\#243](https://github.com/appscode/searchlight/pull/243) ([sajibcse68](https://github.com/sajibcse68))
- Fix version 4.0.0 [\#242](https://github.com/appscode/searchlight/pull/242) ([sajibcse68](https://github.com/sajibcse68))
- Add front matter for docs 4.0.0 [\#241](https://github.com/appscode/searchlight/pull/241) ([sajibcse68](https://github.com/sajibcse68))
- Move alerts under tutorials folder [\#240](https://github.com/appscode/searchlight/pull/240) ([tamalsaha](https://github.com/tamalsaha))
- Make chart namespaced [\#236](https://github.com/appscode/searchlight/pull/236) ([tamalsaha](https://github.com/tamalsaha))
- Change `k8s.io/api/core/v1` pkg alias to core [\#234](https://github.com/appscode/searchlight/pull/234) ([tamalsaha](https://github.com/tamalsaha))
- Use client-go 5.x [\#233](https://github.com/appscode/searchlight/pull/233) ([tamalsaha](https://github.com/tamalsaha))
- Add CustomResourceDefinition [\#232](https://github.com/appscode/searchlight/pull/232) ([tamalsaha](https://github.com/tamalsaha))
- Document how to use kubectl [\#229](https://github.com/appscode/searchlight/pull/229) ([tamalsaha](https://github.com/tamalsaha))
- Add short names for alert objects [\#228](https://github.com/appscode/searchlight/pull/228) ([tamalsaha](https://github.com/tamalsaha))
- Move util to client package [\#227](https://github.com/appscode/searchlight/pull/227) ([tamalsaha](https://github.com/tamalsaha))
- Generate ugorji stuff [\#226](https://github.com/appscode/searchlight/pull/226) ([tamalsaha](https://github.com/tamalsaha))

## [4.0.0](https://github.com/appscode/searchlight/tree/4.0.0) (2017-09-26)
[Full Changelog](https://github.com/appscode/searchlight/compare/4.0.0-rc.0...4.0.0)

**Closed issues:**

- Switch to CustomResourceDefinitions [\#86](https://github.com/appscode/searchlight/issues/86)

**Merged pull requests:**

- Update docs for 4.0.0 release [\#225](https://github.com/appscode/searchlight/pull/225) ([tamalsaha](https://github.com/tamalsaha))
- Install Searchlight as a critical addon [\#224](https://github.com/appscode/searchlight/pull/224) ([tamalsaha](https://github.com/tamalsaha))
- Update chart to add roles for CRD [\#223](https://github.com/appscode/searchlight/pull/223) ([tamalsaha](https://github.com/tamalsaha))
- Revendor errors, log package. [\#222](https://github.com/appscode/searchlight/pull/222) ([tamalsaha](https://github.com/tamalsaha))
- Fixed e2e test [\#221](https://github.com/appscode/searchlight/pull/221) ([aerokite](https://github.com/aerokite))
- Update notifications email templates [\#209](https://github.com/appscode/searchlight/pull/209) ([rubel90](https://github.com/rubel90))

## [4.0.0-rc.0](https://github.com/appscode/searchlight/tree/4.0.0-rc.0) (2017-09-19)
[Full Changelog](https://github.com/appscode/searchlight/compare/3.0.1...4.0.0-rc.0)

**Merged pull requests:**

- Prepare docs for 4.0.0-rc.0 [\#220](https://github.com/appscode/searchlight/pull/220) ([tamalsaha](https://github.com/tamalsaha))
- Update chart to latest convention [\#219](https://github.com/appscode/searchlight/pull/219) ([tamalsaha](https://github.com/tamalsaha))
- Use ObjectReference to write events [\#217](https://github.com/appscode/searchlight/pull/217) ([tamalsaha](https://github.com/tamalsaha))
- Use kubernetes/code-generator [\#215](https://github.com/appscode/searchlight/pull/215) ([tamalsaha](https://github.com/tamalsaha))
- Move all types to types.go [\#214](https://github.com/appscode/searchlight/pull/214) ([tamalsaha](https://github.com/tamalsaha))
- Move analytics collector to root command [\#212](https://github.com/appscode/searchlight/pull/212) ([tamalsaha](https://github.com/tamalsaha))
- Support migration from TPR to CRD [\#211](https://github.com/appscode/searchlight/pull/211) ([aerokite](https://github.com/aerokite))
- Check for ResourceType [\#210](https://github.com/appscode/searchlight/pull/210) ([aerokite](https://github.com/aerokite))
- Use kutil in e2e-test [\#201](https://github.com/appscode/searchlight/pull/201) ([aerokite](https://github.com/aerokite))

## [3.0.1](https://github.com/appscode/searchlight/tree/3.0.1) (2017-08-23)
[Full Changelog](https://github.com/appscode/searchlight/compare/3.0.0...3.0.1)

**Merged pull requests:**

- Find notificaqtion secret from alert namespace [\#207](https://github.com/appscode/searchlight/pull/207) ([tamalsaha](https://github.com/tamalsaha))
- Prepare docs for 3.0.1 release [\#206](https://github.com/appscode/searchlight/pull/206) ([tamalsaha](https://github.com/tamalsaha))
- Update notifier library [\#204](https://github.com/appscode/searchlight/pull/204) ([tamalsaha](https://github.com/tamalsaha))
- Correctly detect pod status [\#203](https://github.com/appscode/searchlight/pull/203) ([tamalsaha](https://github.com/tamalsaha))
- Support patch [\#200](https://github.com/appscode/searchlight/pull/200) ([aerokite](https://github.com/aerokite))
- Example of Alerts for Influx query [\#192](https://github.com/appscode/searchlight/pull/192) ([aerokite](https://github.com/aerokite))

## [3.0.0](https://github.com/appscode/searchlight/tree/3.0.0) (2017-08-21)
[Full Changelog](https://github.com/appscode/searchlight/compare/1.5.9...3.0.0)

**Implemented enhancements:**

- Upgrade Icinga and IcingaWeb2 [\#41](https://github.com/appscode/searchlight/issues/41)

**Closed issues:**

- check\_component\_status - allow specifying component [\#148](https://github.com/appscode/searchlight/issues/148)
- Convert tests to Ginkgo [\#101](https://github.com/appscode/searchlight/issues/101)
- Design page for slack bot\_token info [\#37](https://github.com/appscode/searchlight/issues/37)
- Fix notification messages [\#189](https://github.com/appscode/searchlight/issues/189)
- Document  how to ACK [\#186](https://github.com/appscode/searchlight/issues/186)
- Updating vars does not update in icinga [\#162](https://github.com/appscode/searchlight/issues/162)
- Add warning event if missing notifierSecret [\#158](https://github.com/appscode/searchlight/issues/158)
- Uniform name for check\_\*\*\_volume [\#149](https://github.com/appscode/searchlight/issues/149)
- check\_kube\_event - allow specifying involved object [\#147](https://github.com/appscode/searchlight/issues/147)
- Add Health check for icinga container [\#123](https://github.com/appscode/searchlight/issues/123)
- Fix alert commands [\#117](https://github.com/appscode/searchlight/issues/117)
- Notify cluster admin about soon to be expired certs [\#116](https://github.com/appscode/searchlight/issues/116)
- ENV [\#112](https://github.com/appscode/searchlight/issues/112)
- Handle non-responsive icinga [\#105](https://github.com/appscode/searchlight/issues/105)
- Change User "appscode\_user" to "searchlight\_receiver" [\#103](https://github.com/appscode/searchlight/issues/103)
- Test updated controller manually [\#102](https://github.com/appscode/searchlight/issues/102)
- Fix plugin [\#99](https://github.com/appscode/searchlight/issues/99)
- Cleanup Status of Alerts [\#97](https://github.com/appscode/searchlight/issues/97)
- Use HTTP endpoint for alert ack [\#95](https://github.com/appscode/searchlight/issues/95)
- Use unified notifier [\#93](https://github.com/appscode/searchlight/issues/93)
- Support multiple receivers for each state [\#92](https://github.com/appscode/searchlight/issues/92)
- Support field selectors [\#91](https://github.com/appscode/searchlight/issues/91)
- Use unified notifiers. [\#89](https://github.com/appscode/searchlight/issues/89)
- Automate CA cert generation process [\#84](https://github.com/appscode/searchlight/issues/84)
- Fix secret namespace [\#83](https://github.com/appscode/searchlight/issues/83)
- Fix RBAC [\#82](https://github.com/appscode/searchlight/issues/82)
- Support preferred api group kinds [\#78](https://github.com/appscode/searchlight/issues/78)
- Turn alert.appscode.com/objectType annotations into selectors [\#76](https://github.com/appscode/searchlight/issues/76)
- Cleanup documentation [\#54](https://github.com/appscode/searchlight/issues/54)
- Add tests for check\_volume [\#36](https://github.com/appscode/searchlight/issues/36)

**Merged pull requests:**

- Fix docs [\#199](https://github.com/appscode/searchlight/pull/199) ([tamalsaha](https://github.com/tamalsaha))
- Update package.json version [\#197](https://github.com/appscode/searchlight/pull/197) ([tamalsaha](https://github.com/tamalsaha))
- Make notification email subject informative [\#191](https://github.com/appscode/searchlight/pull/191) ([tamalsaha](https://github.com/tamalsaha))
- Fix README [\#190](https://github.com/appscode/searchlight/pull/190) ([tamalsaha](https://github.com/tamalsaha))
- Fix links [\#185](https://github.com/appscode/searchlight/pull/185) ([tamalsaha](https://github.com/tamalsaha))
- Fix json\_path [\#181](https://github.com/appscode/searchlight/pull/181) ([tamalsaha](https://github.com/tamalsaha))
- Document cluster checks [\#180](https://github.com/appscode/searchlight/pull/180) ([tamalsaha](https://github.com/tamalsaha))
- Document pod & node commands [\#179](https://github.com/appscode/searchlight/pull/179) ([tamalsaha](https://github.com/tamalsaha))
- Make hostfacts listen on all ips by default [\#178](https://github.com/appscode/searchlight/pull/178) ([tamalsaha](https://github.com/tamalsaha))
- Make check\_volume work with minikube [\#177](https://github.com/appscode/searchlight/pull/177) ([tamalsaha](https://github.com/tamalsaha))
- Expose default GO flags [\#176](https://github.com/appscode/searchlight/pull/176) ([tamalsaha](https://github.com/tamalsaha))
- Use go text template for Influx queries. [\#175](https://github.com/appscode/searchlight/pull/175) ([tamalsaha](https://github.com/tamalsaha))
- Fix command docs [\#174](https://github.com/appscode/searchlight/pull/174) ([tamalsaha](https://github.com/tamalsaha))
- Use lowerCamelCase with vars [\#171](https://github.com/appscode/searchlight/pull/171) ([tamalsaha](https://github.com/tamalsaha))
- Fix installer guide [\#170](https://github.com/appscode/searchlight/pull/170) ([tamalsaha](https://github.com/tamalsaha))
- Add images for node & pod alert [\#168](https://github.com/appscode/searchlight/pull/168) ([tamalsaha](https://github.com/tamalsaha))
- Add image for json\_path [\#167](https://github.com/appscode/searchlight/pull/167) ([tamalsaha](https://github.com/tamalsaha))
- Fix Json path plugin [\#165](https://github.com/appscode/searchlight/pull/165) ([tamalsaha](https://github.com/tamalsaha))
- Update event selector [\#164](https://github.com/appscode/searchlight/pull/164) ([tamalsaha](https://github.com/tamalsaha))
- Fix event commands [\#163](https://github.com/appscode/searchlight/pull/163) ([tamalsaha](https://github.com/tamalsaha))
- Fix check commands in icinga image [\#161](https://github.com/appscode/searchlight/pull/161) ([tamalsaha](https://github.com/tamalsaha))
- User docs - part 4 [\#160](https://github.com/appscode/searchlight/pull/160) ([tamalsaha](https://github.com/tamalsaha))
- Check notifier settings [\#157](https://github.com/appscode/searchlight/pull/157) ([tamalsaha](https://github.com/tamalsaha))
- Split data files [\#156](https://github.com/appscode/searchlight/pull/156) ([tamalsaha](https://github.com/tamalsaha))
- User docs - part 3 [\#155](https://github.com/appscode/searchlight/pull/155) ([tamalsaha](https://github.com/tamalsaha))
- User docs - part 2 [\#154](https://github.com/appscode/searchlight/pull/154) ([tamalsaha](https://github.com/tamalsaha))
- Disable command [\#153](https://github.com/appscode/searchlight/pull/153) ([tamalsaha](https://github.com/tamalsaha))
- Update cluster alert plugins [\#152](https://github.com/appscode/searchlight/pull/152) ([tamalsaha](https://github.com/tamalsaha))
- Use uniform name for check\_\*\*\_volume [\#150](https://github.com/appscode/searchlight/pull/150) ([tamalsaha](https://github.com/tamalsaha))
- Server icingaweb2 from root path [\#146](https://github.com/appscode/searchlight/pull/146) ([tamalsaha](https://github.com/tamalsaha))
- Server icingaweb2 from root path [\#145](https://github.com/appscode/searchlight/pull/145) ([tamalsaha](https://github.com/tamalsaha))
- Use cobra with hostfacts [\#144](https://github.com/appscode/searchlight/pull/144) ([tamalsaha](https://github.com/tamalsaha))
- Stop rebuilding postgres during release process. [\#143](https://github.com/appscode/searchlight/pull/143) ([tamalsaha](https://github.com/tamalsaha))
- User docs - part 1 [\#141](https://github.com/appscode/searchlight/pull/141) ([tamalsaha](https://github.com/tamalsaha))
- Fix Docker image issues [\#138](https://github.com/appscode/searchlight/pull/138) ([tamalsaha](https://github.com/tamalsaha))
- Load icinga certs from secret [\#137](https://github.com/appscode/searchlight/pull/137) ([tamalsaha](https://github.com/tamalsaha))
- Fix chart [\#136](https://github.com/appscode/searchlight/pull/136) ([tamalsaha](https://github.com/tamalsaha))
- Fix RBAC permissions [\#135](https://github.com/appscode/searchlight/pull/135) ([tamalsaha](https://github.com/tamalsaha))
- Fix docker images [\#134](https://github.com/appscode/searchlight/pull/134) ([tamalsaha](https://github.com/tamalsaha))
- Enable IcingaWeb2 in alpine image [\#133](https://github.com/appscode/searchlight/pull/133) ([tamalsaha](https://github.com/tamalsaha))
- Fix alpine icinga2 image [\#132](https://github.com/appscode/searchlight/pull/132) ([tamalsaha](https://github.com/tamalsaha))
- Remove AlertStatus [\#131](https://github.com/appscode/searchlight/pull/131) ([tamalsaha](https://github.com/tamalsaha))
- Add livenessProbe for icinga container [\#128](https://github.com/appscode/searchlight/pull/128) ([aerokite](https://github.com/aerokite))
- Add "check\_certificate\_expiry" plugin [\#127](https://github.com/appscode/searchlight/pull/127) ([aerokite](https://github.com/aerokite))
- Bad format except length of 2 & 3 [\#119](https://github.com/appscode/searchlight/pull/119) ([aerokite](https://github.com/aerokite))
- Fix plugins [\#114](https://github.com/appscode/searchlight/pull/114) ([aerokite](https://github.com/aerokite))
- Part 3 - User Guide [\#88](https://github.com/appscode/searchlight/pull/88) ([tamalsaha](https://github.com/tamalsaha))
- Part 1 - User Guide [\#85](https://github.com/appscode/searchlight/pull/85) ([tamalsaha](https://github.com/tamalsaha))
- Various cleanup to chart & deploy script [\#81](https://github.com/appscode/searchlight/pull/81) ([tamalsaha](https://github.com/tamalsaha))
- Use client-go [\#77](https://github.com/appscode/searchlight/pull/77) ([tamalsaha](https://github.com/tamalsaha))
- Enable notification feature in Icinga2 [\#188](https://github.com/appscode/searchlight/pull/188) ([aerokite](https://github.com/aerokite))
- Add concept docs for Alert types [\#173](https://github.com/appscode/searchlight/pull/173) ([tamalsaha](https://github.com/tamalsaha))
- Move notifier secret name inside \*\*\*Alert [\#142](https://github.com/appscode/searchlight/pull/142) ([tamalsaha](https://github.com/tamalsaha))
- e2e test for searchlight [\#120](https://github.com/appscode/searchlight/pull/120) ([aerokite](https://github.com/aerokite))
- Update notification & acknowledgement process [\#115](https://github.com/appscode/searchlight/pull/115) ([aerokite](https://github.com/aerokite))
- Fix jessie image [\#113](https://github.com/appscode/searchlight/pull/113) ([tamalsaha](https://github.com/tamalsaha))
- Add ParseHost\(\) [\#111](https://github.com/appscode/searchlight/pull/111) ([tamalsaha](https://github.com/tamalsaha))
- Prepare alpine based icinga image [\#110](https://github.com/appscode/searchlight/pull/110) ([tamalsaha](https://github.com/tamalsaha))
- Reorganize cmds [\#108](https://github.com/appscode/searchlight/pull/108) ([tamalsaha](https://github.com/tamalsaha))
- Simplify bootstrap process [\#107](https://github.com/appscode/searchlight/pull/107) ([tamalsaha](https://github.com/tamalsaha))
- Fix range loop pointer bugs. [\#106](https://github.com/appscode/searchlight/pull/106) ([tamalsaha](https://github.com/tamalsaha))
- Fix pointer related bugs [\#104](https://github.com/appscode/searchlight/pull/104) ([aerokite](https://github.com/aerokite))
- Reorganize clients [\#96](https://github.com/appscode/searchlight/pull/96) ([tamalsaha](https://github.com/tamalsaha))
- Use separate TPRs for Cluster/Node/Pod alerts [\#90](https://github.com/appscode/searchlight/pull/90) ([tamalsaha](https://github.com/tamalsaha))
- Part 2 - User Guide [\#87](https://github.com/appscode/searchlight/pull/87) ([tamalsaha](https://github.com/tamalsaha))

## [1.5.9](https://github.com/appscode/searchlight/tree/1.5.9) (2017-06-13)
[Full Changelog](https://github.com/appscode/searchlight/compare/1.5.8...1.5.9)

**Closed issues:**

- Add analytics [\#70](https://github.com/appscode/searchlight/issues/70)

**Merged pull requests:**

- Add analytics [\#75](https://github.com/appscode/searchlight/pull/75) ([aerokite](https://github.com/aerokite))
- Explain versioning policy [\#73](https://github.com/appscode/searchlight/pull/73) ([tamalsaha](https://github.com/tamalsaha))
- Change api group to monitoring.appscode.com [\#69](https://github.com/appscode/searchlight/pull/69) ([tamalsaha](https://github.com/tamalsaha))
- Various cleanup for searchlight operator [\#71](https://github.com/appscode/searchlight/pull/71) ([tamalsaha](https://github.com/tamalsaha))
- Add prometheus metrics for hostfacts [\#68](https://github.com/appscode/searchlight/pull/68) ([tamalsaha](https://github.com/tamalsaha))
- Use alpine as the base image for operator [\#67](https://github.com/appscode/searchlight/pull/67) ([tamalsaha](https://github.com/tamalsaha))

## [1.5.8](https://github.com/appscode/searchlight/tree/1.5.8) (2017-05-16)
[Full Changelog](https://github.com/appscode/searchlight/compare/1.5.7...1.5.8)

**Merged pull requests:**

- Use appscode/errors v2 [\#66](https://github.com/appscode/searchlight/pull/66) ([tamalsaha](https://github.com/tamalsaha))

## [1.5.7](https://github.com/appscode/searchlight/tree/1.5.7) (2017-05-03)
[Full Changelog](https://github.com/appscode/searchlight/compare/1.5.6...1.5.7)

**Merged pull requests:**

- Fix for misspell & go\_vet [\#65](https://github.com/appscode/searchlight/pull/65) ([aerokite](https://github.com/aerokite))
- Use updated status types [\#64](https://github.com/appscode/searchlight/pull/64) ([tamalsaha](https://github.com/tamalsaha))
- Fix bugs [\#63](https://github.com/appscode/searchlight/pull/63) ([aerokite](https://github.com/aerokite))
- Run gofmt -s on test pkg [\#62](https://github.com/appscode/searchlight/pull/62) ([tamalsaha](https://github.com/tamalsaha))
- Update docs to new chart location [\#61](https://github.com/appscode/searchlight/pull/61) ([tamalsaha](https://github.com/tamalsaha))
- Move /chart to root directory [\#60](https://github.com/appscode/searchlight/pull/60) ([tamalsaha](https://github.com/tamalsaha))

## [1.5.6](https://github.com/appscode/searchlight/tree/1.5.6) (2017-04-21)
[Full Changelog](https://github.com/appscode/searchlight/compare/1.5.5...1.5.6)

**Implemented enhancements:**

- AlertStatus should keep track of last applied time [\#47](https://github.com/appscode/searchlight/issues/47)

**Merged pull requests:**

- Prepare docs for 1.5.6 release. [\#59](https://github.com/appscode/searchlight/pull/59) ([tamalsaha](https://github.com/tamalsaha))
- Doc fix [\#58](https://github.com/appscode/searchlight/pull/58) ([saumanbiswas](https://github.com/saumanbiswas))
- Stable chart fix [\#57](https://github.com/appscode/searchlight/pull/57) ([saumanbiswas](https://github.com/saumanbiswas))
- Various refinements to chart [\#56](https://github.com/appscode/searchlight/pull/56) ([saumanbiswas](https://github.com/saumanbiswas))
- Update timing fields. [\#55](https://github.com/appscode/searchlight/pull/55) ([tamalsaha](https://github.com/tamalsaha))
- Fixed status fields [\#53](https://github.com/appscode/searchlight/pull/53) ([aerokite](https://github.com/aerokite))

## [1.5.5](https://github.com/appscode/searchlight/tree/1.5.5) (2017-04-19)
[Full Changelog](https://github.com/appscode/searchlight/compare/1.5.4...1.5.5)

**Implemented enhancements:**

-  Modify ports & fix typos [\#50](https://github.com/appscode/searchlight/pull/50) ([aerokite](https://github.com/aerokite))

**Fixed bugs:**

- Fix check\_volume plugin [\#51](https://github.com/appscode/searchlight/pull/51) ([aerokite](https://github.com/aerokite))
-  Modify ports & fix typos [\#50](https://github.com/appscode/searchlight/pull/50) ([aerokite](https://github.com/aerokite))

**Merged pull requests:**

- Add AlertStatus support [\#52](https://github.com/appscode/searchlight/pull/52) ([aerokite](https://github.com/aerokite))

## [1.5.4](https://github.com/appscode/searchlight/tree/1.5.4) (2017-04-16)
[Full Changelog](https://github.com/appscode/searchlight/compare/1.5.3...1.5.4)

**Merged pull requests:**

- Searchlight chart [\#49](https://github.com/appscode/searchlight/pull/49) ([saumanbiswas](https://github.com/saumanbiswas))
- Use thread-safe notifiers. [\#48](https://github.com/appscode/searchlight/pull/48) ([tamalsaha](https://github.com/tamalsaha))

## [1.5.3](https://github.com/appscode/searchlight/tree/1.5.3) (2017-03-01)
[Full Changelog](https://github.com/appscode/searchlight/compare/1.5.2...1.5.3)

**Merged pull requests:**

- Add doc explaining release process. [\#45](https://github.com/appscode/searchlight/pull/45) ([tamalsaha](https://github.com/tamalsaha))
- Add plivo as notifier [\#44](https://github.com/appscode/searchlight/pull/44) ([aerokite](https://github.com/aerokite))

## [1.5.2](https://github.com/appscode/searchlight/tree/1.5.2) (2017-02-27)
[Full Changelog](https://github.com/appscode/searchlight/compare/1.5.1...1.5.2)

**Implemented enhancements:**

- Add slack support as notifier [\#31](https://github.com/appscode/searchlight/issues/31)
- Updated README to highlight supported notifier [\#42](https://github.com/appscode/searchlight/pull/42) ([aerokite](https://github.com/aerokite))
- Added slack support as notifier [\#33](https://github.com/appscode/searchlight/pull/33) ([aerokite](https://github.com/aerokite))
- Unit test [\#30](https://github.com/appscode/searchlight/pull/30) ([aerokite](https://github.com/aerokite))

**Fixed bugs:**

- Secret namespace clarification [\#34](https://github.com/appscode/searchlight/issues/34)
- Use separate flag for namespace [\#38](https://github.com/appscode/searchlight/pull/38) ([aerokite](https://github.com/aerokite))

**Merged pull requests:**

- Update deployment scripts to use the release tag [\#43](https://github.com/appscode/searchlight/pull/43) ([tamalsaha](https://github.com/tamalsaha))
- Various fixes in plugins [\#35](https://github.com/appscode/searchlight/pull/35) ([aerokite](https://github.com/aerokite))

## [1.5.1](https://github.com/appscode/searchlight/tree/1.5.1) (2017-02-14)
[Full Changelog](https://github.com/appscode/searchlight/compare/1.5.0...1.5.1)

**Merged pull requests:**

- While starting Controller, ensure ThirdPartyResource "Alert" [\#29](https://github.com/appscode/searchlight/pull/29) ([aerokite](https://github.com/aerokite))
- Fixed bugs [\#28](https://github.com/appscode/searchlight/pull/28) ([aerokite](https://github.com/aerokite))

## [1.5.0](https://github.com/appscode/searchlight/tree/1.5.0) (2017-02-10)
**Implemented enhancements:**

- Added documentation for parameterized query variable [\#25](https://github.com/appscode/searchlight/pull/25) ([aerokite](https://github.com/aerokite))
- Added unit test for parameterized variables [\#24](https://github.com/appscode/searchlight/pull/24) ([aerokite](https://github.com/aerokite))
- Added authentication for secure hostfacts server [\#22](https://github.com/appscode/searchlight/pull/22) ([aerokite](https://github.com/aerokite))
- Added E2E Tests [\#17](https://github.com/appscode/searchlight/pull/17) ([aerokite](https://github.com/aerokite))
- Added E2E Tests [\#16](https://github.com/appscode/searchlight/pull/16) ([aerokite](https://github.com/aerokite))
- Modified Tests for Icinga2 Custom Plugins [\#13](https://github.com/appscode/searchlight/pull/13) ([aerokite](https://github.com/aerokite))
- Added hostfacts deployment guide [\#6](https://github.com/appscode/searchlight/pull/6) ([aerokite](https://github.com/aerokite))
- Added documentation for notifier [\#3](https://github.com/appscode/searchlight/pull/3) ([aerokite](https://github.com/aerokite))
- Added script to deploy Icinga2 [\#2](https://github.com/appscode/searchlight/pull/2) ([aerokite](https://github.com/aerokite))
- Modified documentation structure [\#1](https://github.com/appscode/searchlight/pull/1) ([aerokite](https://github.com/aerokite))

**Fixed bugs:**

- Added unit test for parameterized variables [\#24](https://github.com/appscode/searchlight/pull/24) ([aerokite](https://github.com/aerokite))
- Replaced petsets support with satefulsets [\#20](https://github.com/appscode/searchlight/pull/20) ([aerokite](https://github.com/aerokite))
- Allow applying alerts while recreating pod with same name [\#18](https://github.com/appscode/searchlight/pull/18) ([aerokite](https://github.com/aerokite))
- Fix Controller [\#15](https://github.com/appscode/searchlight/pull/15) ([aerokite](https://github.com/aerokite))
- Fixing bugs in controller [\#5](https://github.com/appscode/searchlight/pull/5) ([aerokite](https://github.com/aerokite))

**Merged pull requests:**

- Reduced image size [\#27](https://github.com/appscode/searchlight/pull/27) ([aerokite](https://github.com/aerokite))
- Changed Searchlight Controller architectural images [\#26](https://github.com/appscode/searchlight/pull/26) ([aerokite](https://github.com/aerokite))
- Used "appscode/go/net/httpclient" as Client [\#23](https://github.com/appscode/searchlight/pull/23) ([aerokite](https://github.com/aerokite))
- Added E2E Tests [\#21](https://github.com/appscode/searchlight/pull/21) ([aerokite](https://github.com/aerokite))
- Added new EventReason "NoIcingaObjectCreated" for NotFound error [\#19](https://github.com/appscode/searchlight/pull/19) ([aerokite](https://github.com/aerokite))
- Revendor [\#14](https://github.com/appscode/searchlight/pull/14) ([aerokite](https://github.com/aerokite))
- Revendor [\#12](https://github.com/appscode/searchlight/pull/12) ([aerokite](https://github.com/aerokite))
- Revendor [\#11](https://github.com/appscode/searchlight/pull/11) ([aerokite](https://github.com/aerokite))
- modified searchlight controller [\#10](https://github.com/appscode/searchlight/pull/10) ([aerokite](https://github.com/aerokite))
- Used `NamespaceAll` instead of calling multiple API calls [\#8](https://github.com/appscode/searchlight/pull/8) ([aerokite](https://github.com/aerokite))
- Added architectural design guide [\#7](https://github.com/appscode/searchlight/pull/7) ([aerokite](https://github.com/aerokite))
- Detect Icinga2 pod from ancestors [\#4](https://github.com/appscode/searchlight/pull/4) ([aerokite](https://github.com/aerokite))



\* *This Change Log was automatically generated by [github_changelog_generator](https://github.com/skywinder/Github-Changelog-Generator)*