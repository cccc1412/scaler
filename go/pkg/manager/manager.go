/*
Copyright 2023 The Alibaba Cloud Serverless Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package manager

import (
	"fmt"
	"github.com/AliyunContainerService/scaler/go/pkg/config"
	"github.com/AliyunContainerService/scaler/go/pkg/model"
	scaler2 "github.com/AliyunContainerService/scaler/go/pkg/scaler"
	"log"
	"sync"
  "time"
)

var my_map map[string]int

func init_map() {
  my_map = make(map[string]int)
  my_map["d2b844a5da45df270c8d3d5702fb7ef1f8771796"] = 3
  my_map["19353a2e4b72946e52c65c2899dc8597733a3446"] = 3 
  my_map["419076e83e5fb74421897904c58976bb1fa3f9bd"] = 3  
  my_map["ecac08891c8fc642177aefa415abf6ed318f8824"] = 3  
  my_map["4d880ce200606b206f20d2641ff1c892fcc735d3"] = 3  
  my_map["8039ed4fdfbfc943c089ecac39a2cdf5d4c14347"] = 3  
  my_map["bd3456a228222fd75ce1c74098fe7bd3fd1ce59c"] = 3  
  my_map["af97a4ad69b9ca0bf5199fe6f133efb3b7b11368"] = 3  
  my_map["6d4686926c7dd79b48b8586458ce4c4ec0d5e076"] = 3  
  my_map["ac3f183748d38c4d741aee5bf1ca607479fcd557"] = 3  
  my_map["8fd13b83de53b911e072a5e8673a0a29da35c544"] = 3  
  my_map["96bfbddf7a5d83b8db82d3da94f59ce43448c70b"] = 3  
  my_map["324657697c9ded3da4838dc5ab23c71f140f66a9"] = 3  
  my_map["9ff9c48601e28cbabbe0eaeaeff45a880807e7d0"] = 3  
  my_map["15c6be641342abf730d0e71201a8cf854ae18241"] = 3  
  my_map["044bb510661abde83a25447c99dcc2858ea53f05"] = 3  
  my_map["6f8d1ecacb8b8f0a1c2cd05635b01c07bcc0ec3f"] = 3  
  my_map["e3caff0e323e3d0835838826c336d5fc2fb08653"] = 3  
  my_map["e924f541c9495b1c879dfe2bb83a59f74f147a16"] = 3  
  my_map["e3d163254baca649c0675797aadbafd5c592e85e"] = 3  
  my_map["350469b8b8b61e6cc9138b9311668294d699d0d9"] = 3  
  my_map["8371b8baba81aac1ca237e492d7af0d851b4d141"] = 3  
  my_map["9a5548f45c737bd2cf2030fcd3a79bec69a68653"] = 3  
  my_map["e9f83d135d0e7a1b8c72dd3e9adf54e0b4f99b29"] = 3  
  my_map["ba26fa45f55967d5c5e29bb866c85095219db4bd"] = 3  
  my_map["c48b5f84a3d3182a96ad72768ce8159b6ed39afb"] = 3  
  my_map["16ffd629415f1750fc92e3e4a5492d6c7f46763f"] = 3  
  my_map["7fa9eabbe409bd4c5e56edbddf1758a3f677e3b9"] = 3  
  my_map["20980f0efee313898f552d57818c9bbc3f85bc90"] = 3  
  my_map["b49748342822eb2944f0766dadaaf695dcb6f579"] = 3  
  my_map["04b6242af42c84a84e44c9498fd215b10ec9d089"] = 3  
  my_map["a7c19b2360ba9565aff75d8fbee23ff086956f31"] = 3  
  my_map["0135dc8b26775c9bd33154ff279d4440998fa510"] = 3  
  my_map["259d76d439c1b54929d5eb64b9e64792ac780950"] = 3  
  my_map["0ab7d78bff5ebd6034dee1c27c0fd3b1ef4e56ea"] = 3  
  my_map["10c1bce17aab21a8830956a2db4880cc2db1a943"] = 3  
  my_map["251f7aab224feeeb0b58538a68c37f70b7c66d4b"] = 3  
  my_map["35c9b6fcb64e1fde9fedf3b16c7b0fc8536f27af"] = 3  
  my_map["e20ef7c9078598f3505bf588208cfa71832ffe73"] = 3  
  my_map["25a952436dcb76ba56e10e53f556a2a6b78020fd"] = 3  
  my_map["2792e77fd40ef65be920b558674a890545f394f5"] = 3  
  my_map["61a06dbb4f4e0275dbdc18cd7b2ebebeef7ba8ee"] = 3  
  my_map["0fb1c6036216c0e4b98f2930ba6913584a2f9161"] = 3  
  my_map["cf0ce2c6b2a0ba5a5220a99f78e7e72e07189d0e"] = 3  
  my_map["a0b745bfcecce6c5dfa80d8e5302cc2e1ffbf9af"] = 3  
  my_map["ca156df4965406b469db73c518226781035345e8"] = 3  
  my_map["c98e37640ddfac6eee70ffa641bcb8356b4c2219"] = 3  
  my_map["9b3633e36035d8d332691109b38dcc7271cb65c4"] = 3  
  my_map["9104148daf8a2b023bddc3e8fcf06c4daa6c7e03"] = 3  
  my_map["184c73248f46e3e68b9ad0cb5cdd17874db8e8d3"] = 3  
  my_map["f763a4c9f2f06002769c4beb8d380bc3f159df72"] = 3  
  my_map["f31848a1dd75416f01e3486414642a1e8a03239e"] = 3  
  my_map["27cf8a11e075846eb0e695a94020b5151a792e62"] = 3  
  my_map["3c66fdfa3354c35841d82814ec92a826a26f5e24"] = 3  
  my_map["d7a559be16b2f1e8f6dadcc823a66524af5b2441"] = 3  
  my_map["7fa2dad5917eaabfa773b264705447f7f19eb871"] = 3  
  my_map["40a27be5f97a736553f4c6556a276989679109fe"] = 3  
  my_map["2f03d8fe4a87856003dafa9cdd57c77beebad57e"] = 3  
  my_map["be0ad78e7167ef045c8385374b50db5137a5a26f"] = 3  
  my_map["4b831f0e690905f6edce5d9327f4b0d722b445aa"] = 3  
  my_map["8f772af47680677eaa0de758fcc6b241502b0812"] = 3  
  my_map["8ab8fd7f7cfe58aafbee740940425f6ac88487dd"] = 3  
  my_map["206553478c64184fa367506815853592be3279b2"] = 3  
  my_map["2acd0ae94f760719c64933ccd803e0034080c2f4"] = 3  
  my_map["258ce0ee683134d6412b767514771915cef76fde"] = 3  
  my_map["c25fe9996d13d558621a3ac4a4b7c252ef6726a2"] = 3  
  my_map["c34294650ed084e524a6e04841fe5ac95c58fed4"] = 3  
  my_map["3a8e52f3f36179c7a22b1956d0a82fd308c0b644"] = 3  
  my_map["1204111a339bc42f3c921e6f1334424f90bdddc0"] = 3  
  my_map["120a498ff73fe2d84b3f8d64965a53e44f0580e3"] = 3  
  my_map["7b10c5403407184e72f42d137e8760da14fe0640"] = 3  
  my_map["3e9f23288cc3798d2a9a20edef213e8ca617ecb3"] = 3  
  my_map["73d758cccab9b8f9830b1bb2492bca4dc1c9272a"] = 3  
  my_map["1ff2e976d20926a785051733bf188283e79bf8a4"] = 3  
  my_map["93d47fa19d87542357ecd949ae6be3955d95bd23"] = 3  
  my_map["53e1c606590ce266cf04e14dc1be7a960f1c2aa8"] = 3  
  my_map["f68a6fecd7340c73c9debf5e362abf2ef12925ce"] = 3  
  my_map["403f174a20f85c92ebcaa4f83d3581a90d15142d"] = 3  
  my_map["3a417d3e50f502fcab54220b6988227e2ed57475"] = 3  
  my_map["dce8105795da946ad0e927d568f6ad4873f8ad40"] = 3  
  my_map["9ca464d342cdd6545d6b282b6480465b33b10892"] = 3  
  my_map["cc9b4bed0a815816c4014dd8684e808c084dd44b"] = 3  
  my_map["8b1ff8446f0a843e8b114541c15621c1f5e36cd8"] = 3  
  my_map["51e46360734da7f1e1aeed8aff376c4e6faae6b3"] = 3  
  my_map["aba6f081b7f2564d501d250225e2dcd207095fe9"] = 3  
  my_map["64d55a206a567aea06aca78acd4870a60d879c3d"] = 3  
  my_map["bf920007ca148816be1c22b2fc7ad909213ae260"] = 3  
  my_map["eec26d3fb0eab794f87b96d90580d80841eb5238"] = 3  
  my_map["3fa9a52c7dc7f75f6fa5fce85e0e6a8b5c1501ff"] = 3  
  my_map["228a601dd5666dafe50e8dfe728e629f824494eb"] = 3  
  my_map["cb04535115c7723a0fdafd97bc69b25b0aacf6e2"] = 3  
  my_map["d5fbc8a9e6832108e3bd14504386638c63be80d9"] = 3  
  my_map["2c2b37bf286290efcc28d8cac799ad7f1cafd5c2"] = 3  
  my_map["133a934567eda16c8c484b62d11b6006336d7cfb"] = 3  
  my_map["50cd2a2e7ee6a367c3fd011ec4ef394ecb6382c3"] = 3  
  my_map["6da6420c21294ba0d2050cf4134a89f89476140a"] = 3  
  my_map["6f52788ca1e1b70e695d6548f0a76fd3ed044fbe"] = 3  
  my_map["105e1e0e0208ddaaabc59ed6e08aeef437cb1d1a"] = 3  
  my_map["fd2d65a320cdcf3bde0845993b9f6d9095d633cf"] = 3  
  my_map["ba56f1689068ea4f959c2a631c92a898d16d6c69"] = 3  
  my_map["d2b359f422e8be40cab43e7ca1b9c0a9f294540a"] = 3  
  my_map["39684666c1b871093ce43069a022bf51688ac9f6"] = 3  
  my_map["31af38d3c487295ef75f19799f3f15f5f66431b2"] = 3  
  my_map["047d4cfe9f995cc07650a1cc08be97124fe54c4d"] = 3  
  my_map["0badfe300f417844fb302e3d60cc9a80cbd10db6"] = 3  
  my_map["1935a20a0e0b76c423d39cff8c5b42553a5a5f37"] = 3  
  my_map["d256a1a8b60a0f4bb8172a4b3b897e39a4f942fd"] = 3  
  my_map["9a7bdfbe99c82c09e65035f3271faa00b96258b0"] = 3  
  my_map["bda42b5d70c079e20da06c6ddf4f618dc1afe8ef"] = 3  
  my_map["076d225b93ba9442adc0c88325390ca5868ed37e"] = 3  
  my_map["ce0c193237094240a8a81e2e44d961a54472ad45"] = 3  
  my_map["11f05b52bf3157705b1ffde4f6d0f7bd0f3e0dab"] = 3  
  my_map["37b7adc56a88c9e8001576a512878fa5694bca6a"] = 3  
  my_map["2a8a5d163c8167ae0239185de35d53964f7b62db"] = 3  
  my_map["6d8d38f8d7956864f25f428d349a0c94bff3ce14"] = 3  
  my_map["d18e014e3e4d86f4e8c5c8f8837f013b381caa2f"] = 3  
  my_map["cacd029d5b6d301da958a59c236f4c7a7ec36b43"] = 3  
  my_map["ce02ad177a3837abc9a7a06cacd699bf364d3bcc"] = 3  
  my_map["35fb666fe1e5fcc4af1341d2f94038193973f75c"] = 3  
  my_map["e31ac19e3bbb8c0dcf7ceede5fc7b8f6b2f935a6"] = 3  
  my_map["e97e7d5d66536d0e5c854a501951ee4e0104a228"] = 3  
  my_map["4fc78cca3bcf6c3bad5da38e180fdd865a37d1f4"] = 3  
  my_map["3e2042e80001762eafe62df1721db5ec09cdc94e"] = 3  
  my_map["b9eee161b9355279824f7ff806f8b3ebd90adeaa"] = 3  
  my_map["ee85d2b87fa774394e8721e5ad5f8f0aa9be6bbd"] = 3  
  my_map["35f52db5bfc324a3d97c453512617f6cd6c144d3"] = 3  
  my_map["0466f6a37fe342adf5f1787b819c0f51c6f8b306"] = 3  
  my_map["5ae328c76ae5175c98fa1c55d05af959d64defda"] = 3  
  my_map["85bf78f1610dc6a6bddedd05ecc0976d215bae07"] = 3  
  my_map["78fb2956bdc87d380f7f3350520c2b190c52d026"] = 3  
  my_map["7e9ec6b3d661f8dee620d9f12c43560a60076368"] = 3  
  my_map["f008eddd99d2466d2933f65e584fecffb19eefa0"] = 3  
  my_map["1be50bc91ceb9b90183659e8a346005564d546eb"] = 3  
  my_map["a72d5598180b6661659df944d953c1b10e43ac9e"] = 3  
  my_map["db22e5d49c1c4bd38b38adeafddf381f9482bc16"] = 3  
  my_map["bd9b2ddccb8aae15938dc577b6048da1b4085a19"] = 3  
  my_map["bec65708e8640b0a1a2ace81084baad0fc019213"] = 3  
  my_map["62aeb199e061f786cf94541daddaca9c225db708"] = 3  
  my_map["939ee255087c11f667e2d63e99dc727567be33a6"] = 3  
  my_map["88cda92a105df7df419369f53f4b628ba772cacb"] = 3  
  my_map["19f27556a540316ab0d66ddfa84a993c9651b4bb"] = 3  
  my_map["8d2616d3037a928c1437408a5ca5c449a0881a70"] = 3  
  my_map["a8a519733ef9a2ff9cffd8029460b8ebbf713e9d"] = 3  
  my_map["43962874363c4cb35d7c03e841ec7e7fcecf7a5d"] = 3  
  my_map["bcfda646c7117b789ca7a74877eef07e15941bf1"] = 3  
  my_map["51e44dd8be0302e4ef1af18363b1d5ab1b52d094"] = 3  
  my_map["e383f4f75d9f3a561f6bd2c3bd8895aba84083a1"] = 3  
  my_map["6c38ec0453bd2b205ffbf6ca1ac88e488ec112c0"] = 3  
  my_map["dbd1b26cd2044ff081557e16dbe2cd07e42ba449"] = 3  
  my_map["2dca29ea18883b10591d510441f0876f8a6501a7"] = 3  
  my_map["f0bacd392f2be712bfe01e81592226ac8eec59bd"] = 3  
  my_map["7cb003a3fc3a2997a49b04fc0b7416af71ed2410"] = 3  
  my_map["b645a388bea863ef34ea99dbc87612d49d0ad22b"] = 3  
  my_map["ab1a32b43e09d5949cac2301e7858d1055a9303b"] = 3  
  my_map["a38f57ce86ee7712e5d3ea1eb62718b6be1f9c74"] = 3  
  my_map["b242c240278c3826d764faa6dddb77337cd0ac4c"] = 3  
  my_map["a2d9c858cc3d51e607f7403b12488bc0c438e34a"] = 3  
  my_map["7a1ecafc556082b631a12968acf8724dd3bf46e4"] = 3  
  my_map["a559f80e5a6d44a476c2d3410e329ecc550a1c59"] = 3  
  my_map["0dd9b0ae5c58a00ac5bcc9b77100381444a854ad"] = 3  
  my_map["fcd5d154a97f18c1d4c4ef9f1435c2bd668faaf1"] = 3  
  my_map["c19d01f5546bca288a5452f273bef127a1d392e3"] = 3  
  my_map["d6bc9659d81990fb889394a5bda7175de3fc44f8"] = 3  
  my_map["385adb6c2d7c09d85c1845ed0723cd7ef460e0cc"] = 3  
  my_map["c095d9e0c27a664d6902aa2fd6f538ec3ba8894a"] = 3  
  my_map["156e1037e70b69c064acb78a9bc031e415850cc1"] = 3  
  my_map["9543b9e40874abe031f13f750b0be62131e8bd88"] = 3  
  my_map["d8ff5d9917eed0e92970e63fcc7d3da3633aa2a1"] = 3  
  my_map["7b04969d75bbe84bda7ad406eb0397ffe020f251"] = 3  
  my_map["12475a9fcc66e46025459323022765512f36bb5e"] = 3  
  my_map["c58c4f6ff45e2952e462afdfbe388f0bb8a4d665"] = 3  
  my_map["27cb6fa91e401f30bd46ba889cfe73af27096917"] = 3  
  my_map["e19ebb5ea36e72a08c970486ea684f3afa3d72f7"] = 3  
  my_map["d768204d1c2c73e20323a166da6dc8df95a13941"] = 3  
  my_map["092381261f1715c326577cbcd5eb1fe3a9d7412f"] = 3  
  my_map["b529560068b1283e2b2f787184fa001b6d85d2a5"] = 3  
  my_map["ed50293c80d23f6232b677bbff13f9f676893adb"] = 3  
  my_map["e8b63e531a4c0745a4da3b1c4ac2d43d7c2a7eb7"] = 3  
  my_map["870d8fb13c176f5a1a4553890b586280cdefb498"] = 3  
  my_map["4194bbff6ec81d7ff7579279d0787c9496e18a75"] = 3  
  my_map["0e34c8dd3e11af8e111714a5e604f32af2aa4b75"] = 3  
  my_map["030707aacf589606b1160912bdf4396b2d252915"] = 3  
  my_map["2ec38443c8523ed46668ab76ea96f478ca0b9e35"] = 3  
  my_map["58618e7c997d878d50b7c4aab2fb13cc46f12951"] = 3  
  my_map["b11aab406e261a4765e0695494169b5453a77d42"] = 3  
  my_map["c9578f73cf2574f4bed62ef86e14227176b8a1a2"] = 3  
  my_map["7d7462a5d7bd34a6d9a37906d8b05577c6efe277"] = 3  
  my_map["15c0bc0d05e9a8139bc4597f39f1b9cc9853e322"] = 3  
  my_map["3eb8b661855195a7ea5aedc569424f51aa111c15"] = 3  
  my_map["9b1feb62b7d924a90db41192dbdcc0016c6cfea2"] = 3  
  my_map["ef39b277a8a43c4e1efd26dd15c7253cf2fcc5c5"] = 3  
  my_map["2e0396d4dc83c8b69b8ee121c4b96603379f8cb6"] = 3  
  my_map["24c58d8f3b4b85881067a16a02765ba5f7652a07"] = 3  
  my_map["1ce17f6886de49095c034380da8b0629c3942bb8"] = 3  
  my_map["2ca98c95d7fe3ff31beb3c77b0496c4962c18d15"] = 3  
  my_map["fb5c36a212ef3ba2829aa3703c5aa416edc9b318"] = 3  
  my_map["3eef7e4e323166ef66cf80de44936dee0cef1720"] = 3  
  my_map["f0da95d35e3681ae9dfa7be98617b32ec77f8696"] = 3  
  my_map["832b48f7246b6155b3336d9caee05dfb1163ca02"] = 3  
  my_map["ff026b3429d3cabaf43730a3d7064d9568cfaf29"] = 3  
  my_map["02056fa4a980b1d32aecfd15b4702dfe2b1c62e9"] = 3  
  my_map["798a84a12342466b635728285829feb2188e56a6"] = 3  
  my_map["d08146d8e2ed615b44af15557d47893d81664025"] = 3  
  my_map["2d84caf37d0d3e437e97161ea835537c2b4e6c2e"] = 3  
  my_map["6a08b549ac93e5695dc5387837fa2438cdb5a9a1"] = 3  
  my_map["7c91d4c23441ad32c6d1732d82992b30c9d8350c"] = 3  
  my_map["3887a03a71e222b79400ecad2063dd345ba821e6"] = 3  
  my_map["d1a4b760f210b8f3b07cf61847072b7a6c911349"] = 3  
  my_map["3b87beebc8cdf22bbbbf122e85cfc341109815e9"] = 3  
  my_map["938cf56b0a6bc02c808a17fdb62cbbc5c724fd18"] = 3  
  my_map["186de275a231c61091f5018331269db484831044"] = 3  
  my_map["8295bb58ddfe2d6aeaabce0e8a7169c2063618b5"] = 3  
  my_map["dff0fa28d2e42adb2a61f85ce1ec1c08fe3dc7f4"] = 3  
  my_map["ff342e3325c10d2e210319a7461dfd886d80c033"] = 3  
  my_map["516884b51f1eb985077f1aa4d642e21de7d23c49"] = 3  
  my_map["8c0fe750d52bec4d349e5fe14a4b531c49705b10"] = 3  
  my_map["2b0af58326f53b339519e16bc39b761956f7ce33"] = 3  
  my_map["b0d6b73d327f18928d08db5176ccb50170fa771b"] = 3  
  my_map["18972d239f250062cdb2e64084105e72682efe3d"] = 3  
  my_map["16056628af6b93d5a292adf1b01e0d63b01094fc"] = 3  
  my_map["c9b2a885407d79581517313482418e3b54f22e2b"] = 3  
  my_map["24f57fe1f6a354f840032a60f6e2297d583da3b9"] = 3  
  my_map["3829e7509e53c4c05aa8d7858ce0e1995a6c5957"] = 3  
  my_map["1c83fa999efdabe1f64268c77958e44e7951bc97"] = 3  
  my_map["e576a4f2c918c1166dcec40f7958e11c17cc62ad"] = 3  
  my_map["533fd12be07db56503f65abee9e94921b667fc9f"] = 3  
  my_map["3c0af3181f9f4beb377daa45f08bea76db302e5f"] = 3  
  my_map["042ab9eb870d2dc6fc23156ed9cfdb6d025ddf65"] = 3  
  my_map["810152ac99a4960ac19c9b476d711044dd3e143d"] = 3  
  my_map["7b06c1d8353708a17b7fb532b6a92864c638e9fe"] = 3  
  my_map["bf450b5fbe482a04727ed71f391994433f4bdd01"] = 3  
  my_map["2762696c02c752dec332c3ff517f159db5f835cc"] = 3  
  my_map["f9aef1c530ef98ecf6cd1362135599888926df01"] = 3  
  my_map["6d291a81dc709a3681d119937017eb394c68c5a7"] = 3  
  my_map["7c6e126925726e46adc7fd07d609291151fd010a"] = 3  
  my_map["515be8052a95aa44768ce0456439dc80226b8290"] = 3  
  my_map["0b9e2d4f2a934ddecb0c898e828a9807d1c00019"] = 3  
  my_map["64485bf7925722203964938dc334e67fa0e72dcf"] = 3  
  my_map["ced2b0230b020629d32966e523d4898b3774bf8e"] = 3  
  my_map["45a31b0408bf8eb3237c6fb9d8c184ec22bf6ca5"] = 3  
  my_map["cff22db1acaa7264d372573fca4ba337711d66bb"] = 3  
  my_map["56c8c054596faddebfb9c63fc878be276e679dfc"] = 3  
  my_map["437c4d38a98713022a248c909a7c3f8a0a83df38"] = 3  
  my_map["da10c07532a06efb9a6a89baae980c6f1782503b"] = 3  
  my_map["a5d16f04c764b7bd3c5bb36d2ae74e96bf7d126c"] = 3  
  my_map["0181aff53a1108a9365bf812f9707a46a76a93ba"] = 3  
  my_map["7c1846e2a724b584e8a2fcc182a24f47074c7ed6"] = 3  
  my_map["85f071511ca25e5289618e682f3e9676e1a2c039"] = 3  
  my_map["ccc678cae039a57dc53dacd12029967fa838fb36"] = 3  
  my_map["917c4870f1f34e004058bf888c4b0035373ea62b"] = 3  
  my_map["4dd1ea8f772a0a463db394643c5bf517edad8f08"] = 3  
  my_map["c4c2bda8d85d68af1283cd201a59923719c631cf"] = 3  
  my_map["dcc8e7c3c6ff048743330308927fc5a9eee7ef8a"] = 3  
  my_map["3cd0980b4e0f4690491b126fc32b493db7ff127c"] = 3  
  my_map["895c9a3b6fee0ea986cb6bb30af1523a8c69c34d"] = 3  
  my_map["7403a75e31cf84a7752390679a4b66edde70705d"] = 3  
  my_map["e10a70c6996a8fc21e255d91ba695549b46ff02f"] = 3  
  my_map["e14af8cb8708f62fe36bad15283a923d9567471c"] = 3  
  my_map["a22e1d74ccfd1bec69aea09b0a5d681113728e95"] = 3  
  my_map["8caa61e0cc2ae4c14a4273bf640e6630f80c7fd9"] = 3  
  my_map["c38f1d45afbaef7ac379ddeeff8749aa0dab0ccd"] = 3  
  my_map["6c354ee49b7a97c9c213ec55de3dfe4778cfd755"] = 3  
  my_map["ef010ff0357d0bb78ba169ab520ee26fcc6c659d"] = 3  
  my_map["1cb172d056e220c69583c48b7f82201afa7f9165"] = 3  
  my_map["c984081e68aa18b48205f84be24d4219442b395b"] = 3  
  my_map["32ae99416e3e5cd2bd65b4f36581891f30ebfb99"] = 3  
  my_map["e71d14fcc9c51e806cf1001b122c35f775e03d5f"] = 3  
  my_map["4db86de68a7f10e3630763b5a9fbd98fa35ed7e7"] = 3  
  my_map["6c62f7fbf53c40dd68738116f0479eb632285877"] = 3  
  my_map["89b3306fa46997a9e08c8712eb49104eb0897b85"] = 3  
  my_map["9644a25c79b079236312c99090b1b01ccbffddc0"] = 3  
  my_map["2e8773e5d82ab7e1da7f937e2e4f472c36e0bca8"] = 3  
  my_map["67efab166daaeb5455f71289b65845d73e0e03a0"] = 3  
  my_map["4930cabecbd4c5bfcc3c5699ee8a299e70268fa0"] = 3  
  my_map["167e8222a8cb9f19dbc77c280bf0f7fc0d8bb9b6"] = 3  
  my_map["b1b7b1a4933e8590dcaeec4a57442c8d2e16508f"] = 3  
  my_map["796886e0acc5e8e1632bbc5ce788819fd9bcb930"] = 3  
  my_map["368e7269f1810eceb28c4ff43baa0d6413138609"] = 3  
  my_map["e23f1d0bda1ae4dd4fe8e5f4cbf214a883b716f7"] = 3  
  my_map["5cb088b9c161c3ab0e7788da26cc7221f590d373"] = 3  
  my_map["13a319213ec8529389b4bb85c3e829d5f4f62ed8"] = 3  
  my_map["902e0d03249eb85d63d7670ad27bc0992af62eab"] = 3  
  my_map["692b4d2a1fc03b4d2cf0198fcdf0c78cd19d6f07"] = 3  
  my_map["d0621966b40a591c8fa4f905adf3cf1125ab8288"] = 3  
  my_map["c1c54145d7b67bde156ec9ed2db3e93d0554b0df"] = 3  
  my_map["d88307a92209fa938005499792fa17213470fa6b"] = 3  
  my_map["616877714f984e3cf777ee77829c3ebe1304c8a2"] = 3  
  my_map["d0c3b11ee645e3081154243b2516755846b2217c"] = 3  
  my_map["4f003a7743535e73a881723003499dd80ef5c3f3"] = 3  
  my_map["9257f4f46ab366efd08fb0c52125a2f2e93d8024"] = 3  
  my_map["737ca382d836d82c4e2dc2fc32f0f04677c62e67"] = 3  
  my_map["0511ad0ded2003202db6e15a9d0514cad872fd5d"] = 3  
  my_map["b7b41784c6981064affa7b17b92f6d03ffe944bd"] = 3  
  my_map["05c44e9f1e39d84089d252cdc50b08fd2a6bd441"] = 3  
  my_map["02a1fde9f06dbc17079fd3f45fca8d5ed03b8e35"] = 3  
  my_map["85040824d05d037f172aaee016d37c83d502e084"] = 3  
  my_map["c9ca1efec6313976628dc5deb66abe3eeb52ba21"] = 3  
  my_map["f317a510f3454bf597bdc9bf4ab9921a79d2af09"] = 3  
  my_map["696caaced7a8aac11a605aee4c5b7bfa4324a58b"] = 3  
  my_map["45f6c0ab0eb1bee3c4673aec1df8eb025ac741ad"] = 3  
  my_map["6f8d9ac59bfa8eff8c53ddc83ad161b5372f14ef"] = 3  
  my_map["dea060b50591814cb869b9c38f68b6314b2e6f19"] = 3  
  my_map["a313eb3ebcd78877988ed4c68cdfc986e78aa98f"] = 3  
  my_map["1eeb44881951fddba002f6ceac52c324bacc46ce"] = 3  
  my_map["620711a699ef91594017c11ba1c44ae9491a389c"] = 3  
  my_map["85617e958105945d679e7f3d46110f71c611111b"] = 3  
  my_map["695cf76d638a58ab7fad5f3789f923aa1ab5c08b"] = 3  
  my_map["cc0458cbcdb53f351cf22061dd3d9ba3f5d61e37"] = 3  
  my_map["d1b069252d398defdadb07db78d73927e38715a8"] = 3  
  my_map["86017bef6a5eaf982e9e63d774981610549f96ee"] = 3  
  my_map["a5cd5611de6cd9ef20fa3aa0e141cd5062114f4d"] = 3  
  my_map["2efb52e16877d2f3affbd3aecae667416fcade45"] = 3  
  my_map["223231953de5e3a8dd2676d05178364ef8544605"] = 3  
  my_map["87c780770158aba415fdcddc5ed89f262f7e38f2"] = 3  
  my_map["5b6d00d467537593c02dfe735e15c182056c572a"] = 3  
  my_map["2ada289f21360f933228c65bedb30557313b202b"] = 3  
  my_map["4c240fe038c714d8003e4558af0d809a174ce8ba"] = 3  
  my_map["13b6d9e243778f3afa5507aff634810b951174c9"] = 3  
  my_map["8646a0510c930761356308295c24a80c5600a881"] = 3  
  my_map["9183f08a707cd299d4787ff1fae1b3ab40c9cd46"] = 3  
  my_map["3eb351c0c9a9ef2090da4c2bb98e20fd97849e4f"] = 3  
  my_map["7f1c292dd8fbacfd1a80f972f9804305556542f6"] = 3  
  my_map["4fdfa4e850ce267535d1b08f453678c22ef9c4c5"] = 3  
  my_map["14467aa0e9ffcc90342d041b544a2f76a47e7078"] = 3  
  my_map["d23416c29d594a7dade95215db7b1761c633e946"] = 3  
  my_map["8ceaebb4a4b70ff9cb7a2fe79e6551344bee4d20"] = 3  
  my_map["b7624d7917d5dbee72eb189137772afb4dbe1f85"] = 3  
  my_map["547a1c8912fa0073b1bf68f1bf11f647712efd32"] = 3  
  my_map["1eb2b3caed14a490365f1f892c7879e77d51a2b5"] = 3  
  my_map["c5e35a591b9915b413944fa520c63c7d62433d6a"] = 3  
  my_map["f4ae247975f7ea50e851a6e63c6040592f4582f0"] = 3  
  my_map["be5aa35de7a66f708793ea8ee1afd52c87e48ee3"] = 3  
  my_map["fb99b0b135a257b0101bcc179848199be6f8d901"] = 3  
  my_map["4e483ecbd16d436fe88d737b6b5ba3541090c0f9"] = 3  
  my_map["d2b977768081f4d10c52e09c9efc9eaf95a24b66"] = 3  
  my_map["fa31de3e45cbf5aff2152a0d5e426b61c2c1d311"] = 3  
  my_map["d4e5df514f8d652522cf75a32352b80c46d0a794"] = 3  
  my_map["685f1f93218a28889e446d9d61bf2b667011f0f3"] = 3  
  my_map["d3c363b364961902abd9a0b78696201c8224b79e"] = 3  
  my_map["48dd88f7b98458bfcb5b784e740bbf1507a445b6"] = 3  
  my_map["15f33c370204362e346207553e5fd838e8e3a53b"] = 3  
  my_map["8d601befc944a304728043992b5a0a638835956f"] = 3  
  my_map["cea9c9230d765d90eadbcf96f693fde36c00ffb5"] = 3  
  my_map["f0ba6563cdd171aea76993edd4c42484ce788a16"] = 3  
  my_map["ba03a16a233659a9ae2b21923252a3a11ac8a1c4"] = 3  
  my_map["331f4b188d652be9d5db5a08e674a2537c738763"] = 3  
  my_map["a7dc9954f9d331330d181551e8f34d8d240b5f08"] = 3  
  my_map["ade3715d655786af98a2e84c85eb2c379999e695"] = 3  
  my_map["5899841990a05f9c589620937acb9e4353cad31a"] = 3  
  my_map["8c201b1cedc112ef34b68ce53fb2059c201d3c4f"] = 3  
  my_map["78aefca44b00a162ecac9993b8bcb1671f9d7676"] = 3  
  my_map["2f2f2274fd205a673a83bd0ccecf00a08bf8712c"] = 3  
  my_map["99b22214c064c1b52ffbe4f812954007feb28725"] = 3  
  my_map["f7683b06956ceddc5119030e333630f3ec08b31c"] = 3  
  my_map["2629d957eaee009d2da1336402bf3e43dab6d7ea"] = 3  
  my_map["334fe1233fdba26648ff0ecac6b146b32e5f0891"] = 3  
  my_map["7c1c76293ce8ff3d14b4c9dc6484d2250838e7d5"] = 3  
  my_map["615ed4587ec3e1da2ea8eb1b2dae1bd481362918"] = 3  
  my_map["0b330bf5da4adf2552ebf23ce116c28750d8b9c6"] = 3  
  my_map["4e44deea9a60273e060c8ce6e57a08c214cfc509"] = 3  
  my_map["432db71e1c7a682dd3b9e246650eb669bc47eca9"] = 3  
  my_map["d35f206e471cc2f9f7a36f8b0155fa2fa29244ab"] = 3  
  my_map["7351ffbe9e6e81de88ded1644cc6a254dc5b5372"] = 3  
  my_map["6a8f977274dfdc02eb403853f9557649015d0378"] = 3  
  my_map["8816c46e2bed747ec779b29757b26f034733d021"] = 3  
  my_map["e9f5f7ff8346ada131bcf1074252bfe1c52c0e7e"] = 3  
  my_map["1c64defa12c752ef626735a7d3c11003f353abe4"] = 3  
  my_map["729788e1932843953bce53a7d5717bb9cbaf465e"] = 3  
  my_map["7d4c0227f884cb28bffac18fabaf6a588f7a6231"] = 3  
  my_map["7758df61747db33f0af9d8f33d3c985c5465a8be"] = 3  
  my_map["f8163a21fce57cbc4383e8fe42061d6eda96e82c"] = 3  
  my_map["4b7cbec1035f1975e27e1ef7c79f50ad537e791d"] = 3  
  my_map["7c71ab7984213d8b15d01739673e5b8f41699862"] = 3  
  my_map["04e56729fb50644014ce4e9bf06c1038043cdb67"] = 3  
  my_map["d62e9404633e904fae8b8c2478fd4c4e21d63557"] = 3  
  my_map["0b7aa77d990dd331e769a3129c26e4764b688681"] = 3  
  my_map["20952e985be6abf488a639a189491c8d4da36849"] = 3  
  my_map["a10cad3cb5ea178e84ef84f0eb53133a15fb022d"] = 3  
  my_map["a8040abcb04b1322ad4b1e2818e3b1806b463e5e"] = 3  
  my_map["d1cf72ad880dd52e274e96835ff1dbbf1d382673"] = 3  
  my_map["34d0b8062b8b15e344b912ba8190ee9c4f5209b3"] = 3  
  my_map["eef2914d71bca21a699281dc34272fb8079f1fb5"] = 3  
  my_map["6c565c5b6fef64fca92ddf06cb9b7f7fa463d578"] = 3  
  my_map["feb1ad9c80aefd2d55fcaf59fda2ed3e42f6df4e"] = 3  
  my_map["c29671bc3d358047a6ce5711e077e3a18371c644"] = 3  
  my_map["30d511645c1eee8e4abcd2ff15ea5351281546bd"] = 3  
  my_map["51df273902d3ce5c5bc36f216cdd2ac9bf66187f"] = 3  
  my_map["ff0b91b6994e1f28ba9e1bdee09a0381f59ec19c"] = 3  
  my_map["8b83a83f41005c20efd27f7c26a6c7768ede8991"] = 3  
  my_map["12f90b14319e63e4a39eb8eb6f631475a6e1f24b"] = 3  
  my_map["4e7d1d525e7dd090e667ebe8229867f442202eaa"] = 3  
  my_map["ae6e0884a339b774741d7e621290415ebe5f0bf9"] = 3  
  my_map["289e147cdde87c1e1c578efff0323f3f52748f2b"] = 3  
  my_map["e965a0411d841d87645fd53ac1efd6953ef9c292"] = 3  
  my_map["f1e8ef40e67c12241d58b09dfb7349d05413588d"] = 3  
  my_map["b21e5259c41a324c2010b1d187366d448e3584ca"] = 3  
  my_map["837ec95b8aa729911c3d765d8cf897ddf064e2ff"] = 3  
  my_map["0eac8d22a27c6da323865039d81b53cbde0217dd"] = 3  
  my_map["79866642aa723fbee7417bd357d4b737f7983d50"] = 3  
  my_map["3adcc98a3a4283ad9f6ccecca104abd2415f2681"] = 3  
  my_map["14d596f9a66ebf899512758a2eeb9f4200183701"] = 3  
  my_map["3d08a8edbe8f51c67d2f6f5129098e9c4e019577"] = 3  
  my_map["bce5b9d91afcb0b688d1d653a7c64a0913be06f9"] = 3  
  my_map["28a34518177b6713db12b562f1273a50b0a90f39"] = 3  
  my_map["7b63bc5e5f64f462d4e0845713bf99a1154f9599"] = 3  
  my_map["ce66977cc949e05038d2472f846670cc45dc0b4a"] = 3  
  my_map["a07a174131440d912e4ae6ba84675b17f4d63f64"] = 3  
  my_map["fb93c4f77bdabd465b593ad9373cd1881bd6b9d4"] = 3  
  my_map["0b595ae52526e678693d0c18a0eea5092a882edd"] = 3  
  my_map["ae791bf8259a9e5eca4161e55bd7be503794ccec"] = 3  
  my_map["d4641579180eb4c8b59d5a33cc48f3507311aad1"] = 3  
  my_map["e86feaba11f51c6eba53c2eb65c3ca0b2f50f505"] = 3  
  my_map["871b4c09f5b41ebca17c93e27f1ce189e3f1b0ee"] = 3  
  my_map["056480c4241f45fea70f71397d05fd4dde6ca5ad"] = 3  
  my_map["7305a06dc2f183f4962dd4ab1f21e27dcc2b13ba"] = 3  
  my_map["baea84e18d1ee3c310878708d956104fac092feb"] = 3  
  my_map["efb2e10c19b6475a67fe4d8d687c13e8833d5193"] = 3  
  my_map["4783e8abf22d2781ff1f98c21a235be767513df0"] = 3  
  my_map["ac7ddca243c658e0d6751dcfd5c5b4ce3796d948"] = 3  
  my_map["0fa9bc7c26384f1cb7c03d0757f0da3741151d0a"] = 3  
  my_map["bfe1ad9927f9fbcc1b9989d7fc82ca92fabfc6a4"] = 3  
  my_map["108b8602d1c930790105ccf4b072935c9196abce"] = 3  
  my_map["65bf63c09c25baab110fde7d27fef79b50c57fce"] = 3  
  my_map["89d427994f8f92a515d5f104b22696382c1ac242"] = 3  
  my_map["d3a3f36e6aa8caf20215fcb8846c0e3a28192dd2"] = 3  
  my_map["9b69aab26336b30c56a07e14727575f57942d662"] = 3  
  my_map["75366afde1034dd5acb78d435a02da09715382ff"] = 3  
  my_map["2ecfc67d09e13ae95216378800ead137a419798f"] = 3  
  my_map["ec1153bfb8c259a8dc2fea940f30bb479ae40ae4"] = 3  
  my_map["0ca93b5bf30473fd0161cac884e8f00450fb0857"] = 3  
  my_map["15f1f3fbe67063c2246a1234713fe47020d4337f"] = 3  
  my_map["916450d5995575b5067d241bc128ebf36f6ad5e5"] = 3  
  my_map["b5c97889b5d74fe71b41ce60e90994dd6869556c"] = 3  
  my_map["0d65f30b04c2e634d227551016d8cfc1cc6e111a"] = 3  
  my_map["cdb668b2874ba56d764f22a4a797cb25bc6e86b4"] = 3  
  my_map["3e8e4afde47c85462437427b85c85383859d53ec"] = 3  
  my_map["391c0733fec0b1a24d8d6c35b8b68a96a2ed5644"] = 3  
  my_map["0ce048366c22cddd368f17658fc7a638fdb8dbc9"] = 3  
  my_map["203be9fc2b6d16d82911c083fc839c20193fb637"] = 3  
  my_map["dfd9a820c6f3dcaffea1649adc248d7a95bb01ed"] = 3  
  my_map["087b588c7d2c7ba345199981f36896ee90be9cf0"] = 3  
  my_map["cce7a4ecaf6e5c90887396990fe4a4dd1e6118a0"] = 3  
  my_map["85511889bd9d53218fb62efdc26aa10bd1c6519e"] = 3  
  my_map["5a1aa268f52ba7dbaa388d9af847e494b153670a"] = 3  
  my_map["3b99d30f6314969fe08aefdd43400be7bd0fd9ea"] = 3  
  my_map["aa01f5cb0c72e80e14226d9351dfa6bd72e718ca"] = 3  
  my_map["478a845b90e2f829ea7ebe073faf257c6ad4e476"] = 3  
  my_map["6649c17b1b626a527072f73194992c09c390bed6"] = 3  
  my_map["2005f2f3c8d2f57d74f873a993a550d968ec027f"] = 3  
  my_map["a94fe6e097d195587ecbeee1eb005e4d81322c2a"] = 3  
  my_map["173c7786f45f8e128c16f9a0f57871cea8585994"] = 3  
  my_map["ed277f6d99e8062c93da008ee334e0bd2fb1e8fe"] = 3  
  my_map["037ee604ba07d956d1e17f65f74ce2c8a104be66"] = 3  
  my_map["63a65fdc1736fb371656a6f80eb0d0ac6929f767"] = 3  
  my_map["02e6d032c352a70e31e0a9faffebf8b81e363234"] = 3  
  my_map["efa3aa2c495aea4cd675db7cec949fa899b69a7f"] = 3  
  my_map["c4ad822e4ee94b6d8e9daae0ac5899b2b3fd5fbc"] = 3  
  my_map["7f9de1501322b093756c574e787fc8a85768ffdb"] = 3  
  my_map["f24ce71ecf08e7025888b82344a30eb3659db9e4"] = 3  
  my_map["3c6838a4c2e66f2e65c35465adc39806a92f4b71"] = 3  
  my_map["89342ee2994154f441098093f5e07d7972d2a5ee"] = 3  
  my_map["13120be6e45d079313d9568cbf32a1b938ea6d09"] = 3  
  my_map["c5ae61d6d75b1948b71f56d49383c42ffe5ac684"] = 3  
  my_map["2a34a0345ed7d1ce1473ef815a82db7d0d5a9221"] = 3  
  my_map["fbc7ad55f171aef10115623a2c77571110427dc8"] = 3  
  my_map["e29fd93fc8c2531ce5448fcfcd4ba5bcfb2009ee"] = 3  
  my_map["241e17da836a748f390175d3a119a69c859660f5"] = 3  
  my_map["02b0622f917d08394ca4deff84a5ad666e88a26b"] = 3  
  my_map["30ea5d5887fa8843d4543c6349d51b04b1d79dcc"] = 3  
  my_map["fcfd26394dc5016efd438979dfcc9ad069a212fe"] = 3  
  my_map["4505ee1ca189e2eb39f5992d5834440c90a2809b"] = 3  
  my_map["7d94adca9dab87296b08337ee093df214b4c5ccd"] = 3  
  my_map["41b5a92d14470e92c543401dcb4f78c40251a74f"] = 3  
  my_map["460eba776000f852460aa2f24b7c2c2240570683"] = 3  
  my_map["c9dff4b0ba3c7b0eb85849ec91854d2602b75e7e"] = 3  
  my_map["dbbd1b2c40f49b81060affec29c6d6f3b0a58ff2"] = 3  
  my_map["cac3feff9bf6917beb7cde29602e56f6426feff7"] = 3  
  my_map["c01d4ef6985b526e1b325a9dbcfe5f316340923d"] = 3  
  my_map["7f6abc6c82fd4d65bc88b7138d1c70e32d11788f"] = 3  
  my_map["e74a9b5500bfcd5227844173e99105fc3f872b38"] = 3  
  my_map["df206ebd542eaf3c04031f0904f8b9e082e168f6"] = 3  
  my_map["58b775c41e139fd80b788728ecdc2ec20edb370e"] = 3  
  my_map["6010cbf1905b30ad33222937079f839386f956f6"] = 3  
  my_map["f39d48da40c066ac8f090cc8a061016b16b076b2"] = 3  
  my_map["f5913bfd8854f626960c56f248526bf2fac57e2d"] = 3  
  my_map["b53d1c58ed7f09e58f3f636c6adc1fa3aec04e36"] = 3  
  my_map["d39a6ac5257f9a1f2c45d0beabf935b6fc6503ef"] = 3  
  my_map["4715dab6c29fa12e153d299995f65518d36a7391"] = 3  
  my_map["ff1caca254135e26f13db874a752d518bb0cbb00"] = 3  
  my_map["b245e34d0b83ea1e8bf61d588eff65c9c84912d5"] = 3  
  my_map["592fb2b4107e99b10cda4a05707964479dfbf7a7"] = 3  
  my_map["47ef06ca4d638caee1f7273eb4c348bc28f3ff25"] = 3  
  my_map["eeb794c35e3d1a95838556f712e779c11b210785"] = 3  
  my_map["6dc043cd6a1285caf51a60d2e1f78cc31f906383"] = 3  
  my_map["468cf942d2b98337e17fdd77a7f0c020966a65f1"] = 3  
  my_map["73f90984bee008d2523b9551b7fc6e7cd782463c"] = 3  
  my_map["55179f0a18d43fbe4ea513054a2064decddd075c"] = 3  
  my_map["1e7698c31ae6905ab925e03710e425c96009f3f3"] = 3  
  my_map["8d28d92b5b1300f0fcc1a8c194f9bc955d3c72c5"] = 3  
  my_map["d91c4e20d8ef9dcc9a9aa2bf82c2b14d563cd3f0"] = 3  
  my_map["8a0be0de8f9f2974c20da448b4bf9958d2c46a67"] = 3  
  my_map["cff50c60b8d3c973c868828ca0bb45a1d67a1180"] = 3  
  my_map["2b3d197d2a5076d186a9784a9820a18e8071ed09"] = 3  
  my_map["6dd9e76a7815bd46d234268da11c970e1e1ddbed"] = 3  
  my_map["feb1074952b5118acd3cd97615e8a4fabad0960a"] = 3  
  my_map["91f9d7a0d072c0d9efc7f91f99676068db8ca86a"] = 3  
  my_map["cb849261f41a10ae6773442ea753a29509f77105"] = 3  
  my_map["dbf3656d161e6db006b06c2c48ef3fbe02f39963"] = 3  
  my_map["d66c1db424bbe4706bb8a33dc197cf0daa48c42a"] = 3  
  my_map["e1feffe30a5b382d0ff9492aa756257ba4afe59a"] = 3  
  my_map["9cdb9644ccd8e9a99f4d95c7a3e4bdebb181e717"] = 3  
  my_map["1803d6b75561cc0cae4449f6649f7ac901493d8f"] = 3  
  my_map["ea5605ff70625873cd8b28a943c519d73263fac6"] = 3  
  my_map["04bdb20b3c846df42534237a5c4bfe2637b89f77"] = 3  
  my_map["26fa64bd16d545ab31f730622150ee598f8bedf2"] = 3  
  my_map["9131a46a0fdfc1a83ed92483ded76ab669d9a11a"] = 3  
  my_map["85145d4d8cd946b63612f5173dec5893d1f928e0"] = 3  
  my_map["b90ef1147456349e9bd10a67927b1b2a512dfc28"] = 3  
  my_map["5b921017f2baa34c30991ca3be0516b7687858a3"] = 3  
  my_map["33bed9cb5a2a9d6762350d1dfae8898f19254202"] = 3  
  my_map["a212cf9c3981ce7b6f541c47d90783e4fcbd2c1b"] = 3  
  my_map["e08ffb059ceaedb5937e97b6eccd187863ec7e5d"] = 3  
  my_map["bb54bcfee0d586a9606e399ddd93e6e7ba409ffa"] = 3  
  my_map["2bbbe6881a5e5fcb0a36d96422ca231a821e2058"] = 3  
  my_map["00338771693db69a22ef177920cdc8f865bd85a7"] = 3  
  my_map["d42ea7c465a254e5465070c26eb0db4f41a27d88"] = 3  
  my_map["abbcd623ff15abf449b90914bbe4d0d53ca6d64d"] = 3  
  my_map["3e06a372ddf2d4e5c41cf138e5371d0169b7a202"] = 3  
  my_map["44c806a225a75f212957270ca7ec97e36f9de1ca"] = 3  
  my_map["d80e59db402d99886aa697560e23a8fc11787916"] = 3  
  my_map["ee38fba1af2e34d3c0ce43e6d68247ec7078c571"] = 3  
  my_map["6565b1fc5dda7d4a482b17ab57ec6d1d885dbaf3"] = 3  
  my_map["3094edaf6e3e76922ff8a72a90672e1dbb18e7df"] = 3  
  my_map["a06b3a34b8443777baeef4e19e5a80b0d0527df6"] = 3  
  my_map["a0346a0239893eecdadc20e6fee4e000e3c374d1"] = 3  
  my_map["a971cbd8b5ef198a52ba9fea819f5a2166c07131"] = 3  
  my_map["4ea6b320928a7700a011553085d0011b03274bab"] = 3  
  my_map["0a0ca59b77959302611c57d75ffc599485c12e34"] = 3  
  my_map["6f83e25d1fad0b50fe5434423db84731f9a166c2"] = 3  
  my_map["9c1222b7b54433655b7c623555490f47042cafdc"] = 3  
  my_map["691e44923cf182f29b793be363ed27de32cdb274"] = 3  
  my_map["663e9d923bb282b4624f45707582ff2e33666a1b"] = 3  
  my_map["843fd86a7b09f0fe324f8c1a505e915eb8db7e0c"] = 3  
  my_map["35abc2d27b129b752c68e6e45c266194575eebc2"] = 3  
  my_map["9d7398776642f6deae913b889748a38577637077"] = 3  
  my_map["2f2f4172d4a070b1fbf83d449df406e2c5e1cf89"] = 3  
  my_map["f43c6f855fce936da9177b65dfc67ef579edf77e"] = 3  
  my_map["e16cfbfb91fe4a0c13d0b7e186977c751adade62"] = 3  
  my_map["713401e8e946f27b03ec8dd51ddefab71207f20e"] = 3  
  my_map["9a89d9b5b379a958255fd34e5ee770ace2cdafb2"] = 3  
  my_map["13ff97e1aa97e16e72cf7e869e7f7c48ba46a2c4"] = 3  
  my_map["280c68a3dbe9b7936e81d0927ded51437e1cfb20"] = 3  
  my_map["0fc98621700fde5998e7620e9a88688b7d2afbe4"] = 3  
  my_map["f2ce5b954bb55d1c30b47dc12c66cc565961fcd7"] = 3  
  my_map["db697eb65ea626e8a02baeb261e315210891fea8"] = 3  
  my_map["2f052f4e35a26908c03b0d7c0d8780ba806aefc9"] = 3  
  my_map["029698747f84935a1b07bb97a1ffb1d40fc8deca"] = 3  
  my_map["80f59fac3598c93a807eb04b053d28a716f99af3"] = 3  
  my_map["9d1f3cce4ca7d5801a3f961c4acfa440f89e20b1"] = 3  
  my_map["69c752536f36226a7e6e5981b7aee5adee94ff6b"] = 3  
  my_map["97256a15f348b0350a93129691ebdafd14227c9d"] = 3  
  my_map["8382ace4b29af71dccfc9bf259fbf7b82cdd1413"] = 3  
  my_map["79a93e2460cf909bcd27e1a220c11fd482a70cec"] = 3  
  my_map["a22beea5881245fc842dbdec3b7fce85f8bdc12b"] = 3  
  my_map["39476527c2c0688970b71bffd453efd386a5e7f1"] = 3  
  my_map["3757ca10ef8744acbe4ea05fee5b64cbc0e33998"] = 3  
  my_map["a515e0378196b59412af9070160660351cad3553"] = 3  
  my_map["8d153d7f7c7b257039dcfbe71cc9f8fb5de08703"] = 3  
  my_map["ca79da4f9177bf3176db3817cca12189ab19335d"] = 3  
  my_map["adc49a6fc949e5ad693785dd96c86b36f267688e"] = 3  
  my_map["384004eff6c838d25962810585bbce5d71f0a877"] = 3  
  my_map["d686e7e976645b7b5883f658c943264fd3a0686c"] = 3  
  my_map["dd91924ca522df94bfbb19a135268e3ebf7b86d7"] = 3  
  my_map["a7475c9189f0075180926fdad4fe515e90dfeca6"] = 3  
  my_map["1985e23a0856041b1173345c2a7229127f4025e2"] = 3  
  my_map["635388c9afcb28ada1cf9e3edfac8558ab68d7ea"] = 3  
  my_map["3b688bdd12d280a9b7fd3518881a7251753a8b35"] = 3  
  my_map["ecededdec4d51fb740472ba33ba7fb13304edfb9"] = 3  
  my_map["36cc3d741b86ad4298c3cdc279ff07f01298f5d9"] = 3  
  my_map["ba815b764ce640a42ed3fc8e62108086a3527749"] = 3  
  my_map["1cebd037db566e412f504826cb50effaf906d920"] = 3  
  my_map["f11d91a25894dbaf414c00568bbfa6aa4358d83b"] = 3  
  my_map["d8b6cb3dff22588dc012a791951c27880328f4cd"] = 3  
  my_map["2a67e4bc44ee4f68824a33205388f5b84f3883f6"] = 3  
  my_map["91dc0d70cdaf2c8776d67070a5ed1b4fc56d94cc"] = 3  
  my_map["b3ab2d05c50882eeea5d96779cde650558078b22"] = 3  
  my_map["5b578a6b1c1b9450d92041199609aa073868f308"] = 3  
  my_map["eb85e9b8f47351d5a4b5bd805a0c6331ce69f0d5"] = 3  
  my_map["dc257b03db32f877738b4d6646a95f22ca588f13"] = 3  
  my_map["b3cad51f1de6e8e40c112775f5e7812e3d584bdb"] = 3  
  my_map["32e2365ac01be7f166a2b9c8ee9d6c2093002e78"] = 3  
  my_map["b2af3fc6fa60b5d6123cfbe3d3b7041971801544"] = 3  
  my_map["89d061dc51ccc16c0a8756084c61ea8a3e01eb0d"] = 3  
  my_map["d7d06a5e685e4c3465f036c8530e895030cc22e0"] = 3  
  my_map["b335cf32e0346f806e506633cee03ba4094fbca6"] = 3  
  my_map["c5ec647cf1d1fe393d60d02ed01b702e1e4ae545"] = 3  
  my_map["336ed51d1ab1cd7dd184d3362fe62b075b75693f"] = 3  
  my_map["5838c15a9c19bba2d23ccda0620a16728a361281"] = 3  
  my_map["7c446cb6f9ade24af61bdea12bbabdf19bf61ac2"] = 3  
  my_map["09940f8b95ae0351492eccab24247404ad9c1c69"] = 3  
  my_map["a0c7da07974f39ffb7c3aecf29e41bbadb61942f"] = 3  
  my_map["7ba8ce84afa16f55aa4846ba831df1fe534eec39"] = 3  
  my_map["9317abbcf1b34ce2d22ebf45779c0977d963ab33"] = 3  
  my_map["423eebeea3d2d8dd7d4c460c25420961a4338b02"] = 3  
  my_map["ce9d9dee935af179e37766a9c4097f722312bbcf"] = 3  
  my_map["e7379454c4e4fc248e57eebbbead54ad06d41cd2"] = 3  
  my_map["1de3e81af891eaa431ef3ab4d857c5fe46140d6b"] = 3  
  my_map["ce9a3976d666f967e1ae802f74ed51d58a2fee97"] = 3  
  my_map["92ea72217b73f1d711a7dd007def2a73abadecac"] = 3  
  my_map["e48366b1aa5c42a01ff12a6f82df8d847811eb02"] = 3  
  my_map["0374f9c139b24cfea4b4bba2b4ac26bf265e4f4b"] = 3  
  my_map["d4a7c578dfe155db627dcafe87399032514b5318"] = 3  

  my_map["nodes2"] = 2
  my_map["roles2"] = 2
  my_map["rolebindings2"] = 2
  my_map["certificatesigningrequests2"] = 2
  my_map["binding2"] = 2
  my_map["csinodes2"] = 2

  my_map["nodes1"] = 1
  my_map["roles1"] = 4
  my_map["rolebindings1"] = 4
  my_map["certificatesigningrequests1"] = 4
  my_map["binding1"] = 3
  my_map["csinodes1"] = 1

}

type Manager struct {
	rw         sync.RWMutex
	schedulers map[string]scaler2.Scaler
	config     *config.Config
}

func New(config *config.Config) *Manager {
  init_map()
	return &Manager{
		rw:         sync.RWMutex{},
		schedulers: make(map[string]scaler2.Scaler),
		config:     config,
	}
}

func (m *Manager) GetOrCreate(metaData *model.Meta) scaler2.Scaler {
	m.rw.RLock()
	if scheduler := m.schedulers[metaData.Key]; scheduler != nil {
		m.rw.RUnlock()
		return scheduler
	}
	m.rw.RUnlock()

	m.rw.Lock()
	if scheduler := m.schedulers[metaData.Key]; scheduler != nil {
		m.rw.Unlock()
		return scheduler
	}
	log.Printf("Create new scaler for app %s", metaData.Key)
  if(my_map[metaData.Key] == 3) {
    var DefaultConfig3 = config.Config{
	    ClientAddr:           "127.0.0.1:50051",
	    KeepAliveInterval:    15 * 60 * 1000,
      // PreloadInterval:       0, 
      GcDuration: 1 * 1000 * time.Millisecond,
      PreloadDuration: 90 * 1000 * time.Millisecond,
    }
    scheduler := scaler2.New3(metaData, &DefaultConfig3)
    m.schedulers[metaData.Key] = scheduler
	  m.rw.Unlock()
	  return scheduler
  }
  if(my_map[metaData.Key] == 1) { //保留一个，不释放
    var DefaultConfig3 = config.Config{
	    ClientAddr:           "127.0.0.1:50051",
	    KeepAliveInterval:    150 * 60 * 1000,
      // PreloadInterval:       0, 
      GcDuration: 1 * 1000 * time.Millisecond,
      PreloadDuration: 90 * 1000 * time.Millisecond,
    }
    scheduler := scaler2.NewLow(metaData, &DefaultConfig3)
    m.schedulers[metaData.Key] = scheduler
	  m.rw.Unlock()
	  return scheduler
  }
  if(my_map[metaData.Key] == 4) { //最低水位限制，一次load 10 20 30
    var DefaultConfig3 = config.Config{
	    ClientAddr:           "127.0.0.1:50051",
	    KeepAliveInterval:    15 * 60 * 1000,
      // PreloadInterval:       0, 
      GcDuration: 1 * 1000 * time.Millisecond,
      PreloadDuration: 90 * 1000 * time.Millisecond,
    }
    scheduler := scaler2.NewWater(metaData, &DefaultConfig3)
    m.schedulers[metaData.Key] = scheduler
	  m.rw.Unlock()
	  return scheduler
  }

	scheduler := scaler2.New(metaData, m.config)
	m.schedulers[metaData.Key] = scheduler
	m.rw.Unlock()
	return scheduler
}

func (m *Manager) Get(metaKey string) (scaler2.Scaler, error) {
	m.rw.RLock()
	defer m.rw.RUnlock()
	if scheduler := m.schedulers[metaKey]; scheduler != nil {
		return scheduler, nil
	}
	return nil, fmt.Errorf("scaler of app: %s not found", metaKey)
}
