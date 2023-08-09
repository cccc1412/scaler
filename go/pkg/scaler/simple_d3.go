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

package scaler

import (
	"container/list"
	"context"
	"fmt"
	"github.com/AliyunContainerService/scaler/go/pkg/config"
	model2 "github.com/AliyunContainerService/scaler/go/pkg/model"
	platform_client2 "github.com/AliyunContainerService/scaler/go/pkg/platform_client"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"sync"
	"time"

	pb "github.com/AliyunContainerService/scaler/proto"
	"github.com/google/uuid"
)

type SimpleD3 struct {
	config         *config.ConfigD3
	metaData       *model2.Meta
	platformClient platform_client2.Client
	mu             sync.Mutex
	wg             sync.WaitGroup
	instances      map[string]*model2.Instance
	idleInstance   *list.List

  assign_time   int64
}

var init_duration_map map[string]int
func init_dm() {
  init_duration_map["13b6d9e243778f3afa5507aff634810b951174c9"]=1515
  init_duration_map["d256a1a8b60a0f4bb8172a4b3b897e39a4f942fd"]=137
  init_duration_map["384004eff6c838d25962810585bbce5d71f0a877"]=8936
  init_duration_map["810152ac99a4960ac19c9b476d711044dd3e143d"]=207
  init_duration_map["7c1c76293ce8ff3d14b4c9dc6484d2250838e7d5"]=685
  init_duration_map["798a84a12342466b635728285829feb2188e56a6"]=551
  init_duration_map["13ff97e1aa97e16e72cf7e869e7f7c48ba46a2c4"]=8586
  init_duration_map["2f052f4e35a26908c03b0d7c0d8780ba806aefc9"]=9533
  init_duration_map["e74a9b5500bfcd5227844173e99105fc3f872b38"]=182
  init_duration_map["b3ab2d05c50882eeea5d96779cde650558078b22"]=7442
  init_duration_map["0b595ae52526e678693d0c18a0eea5092a882edd"]=172
  init_duration_map["4930cabecbd4c5bfcc3c5699ee8a299e70268fa0"]=1651
  init_duration_map["b2af3fc6fa60b5d6123cfbe3d3b7041971801544"]=1623
  init_duration_map["d8b6cb3dff22588dc012a791951c27880328f4cd"]=5035
  init_duration_map["7305a06dc2f183f4962dd4ab1f21e27dcc2b13ba"]=116
  init_duration_map["696caaced7a8aac11a605aee4c5b7bfa4324a58b"]=1183
  init_duration_map["28a34518177b6713db12b562f1273a50b0a90f39"]=135
  init_duration_map["6a08b549ac93e5695dc5387837fa2438cdb5a9a1"]=171
  init_duration_map["a94fe6e097d195587ecbeee1eb005e4d81322c2a"]=1850
  init_duration_map["fb99b0b135a257b0101bcc179848199be6f8d901"]=362
  init_duration_map["2ec38443c8523ed46668ab76ea96f478ca0b9e35"]=286
  init_duration_map["ff026b3429d3cabaf43730a3d7064d9568cfaf29"]=189
  init_duration_map["5899841990a05f9c589620937acb9e4353cad31a"]=1542
  init_duration_map["cf0ce2c6b2a0ba5a5220a99f78e7e72e07189d0e"]=138
  init_duration_map["184c73248f46e3e68b9ad0cb5cdd17874db8e8d3"]=171
  init_duration_map["cdb668b2874ba56d764f22a4a797cb25bc6e86b4"]=190
  init_duration_map["368e7269f1810eceb28c4ff43baa0d6413138609"]=1478
  init_duration_map["916450d5995575b5067d241bc128ebf36f6ad5e5"]=187
  init_duration_map["e97e7d5d66536d0e5c854a501951ee4e0104a228"]=198
  init_duration_map["f7683b06956ceddc5119030e333630f3ec08b31c"]=655
  init_duration_map["8382ace4b29af71dccfc9bf259fbf7b82cdd1413"]=7829
  init_duration_map["79a93e2460cf909bcd27e1a220c11fd482a70cec"]=5963
  init_duration_map["2005f2f3c8d2f57d74f873a993a550d968ec027f"]=809
  init_duration_map["d4a7c578dfe155db627dcafe87399032514b5318"]=4331
  init_duration_map["d686e7e976645b7b5883f658c943264fd3a0686c"]=1985
  init_duration_map["3b87beebc8cdf22bbbbf122e85cfc341109815e9"]=424
  init_duration_map["b3cad51f1de6e8e40c112775f5e7812e3d584bdb"]=8724
  init_duration_map["80f59fac3598c93a807eb04b053d28a716f99af3"]=2707
  init_duration_map["0fc98621700fde5998e7620e9a88688b7d2afbe4"]=3184
  init_duration_map["9183f08a707cd299d4787ff1fae1b3ab40c9cd46"]=327
  init_duration_map["6d8d38f8d7956864f25f428d349a0c94bff3ce14"]=3880
  init_duration_map["8295bb58ddfe2d6aeaabce0e8a7169c2063618b5"]=142
  init_duration_map["5838c15a9c19bba2d23ccda0620a16728a361281"]=2284
  init_duration_map["e29fd93fc8c2531ce5448fcfcd4ba5bcfb2009ee"]=1137
  init_duration_map["029698747f84935a1b07bb97a1ffb1d40fc8deca"]=7604
  init_duration_map["b11aab406e261a4765e0695494169b5453a77d42"]=144
  init_duration_map["385adb6c2d7c09d85c1845ed0723cd7ef460e0cc"]=814
  init_duration_map["1c64defa12c752ef626735a7d3c11003f353abe4"]=1427
  init_duration_map["9317abbcf1b34ce2d22ebf45779c0977d963ab33"]=7097
  init_duration_map["3f81d2da363ef934f9b96b1f2b99b271be28d847"]=1297
  init_duration_map["0b9e2d4f2a934ddecb0c898e828a9807d1c00019"]=1285
  init_duration_map["8b83a83f41005c20efd27f7c26a6c7768ede8991"]=3679
  init_duration_map["087b588c7d2c7ba345199981f36896ee90be9cf0"]=1205
  init_duration_map["4194bbff6ec81d7ff7579279d0787c9496e18a75"]=153
  init_duration_map["e16cfbfb91fe4a0c13d0b7e186977c751adade62"]=5912
  init_duration_map["bec65708e8640b0a1a2ace81084baad0fc019213"]=170
  init_duration_map["6c565c5b6fef64fca92ddf06cb9b7f7fa463d578"]=406
  init_duration_map["0374f9c139b24cfea4b4bba2b4ac26bf265e4f4b"]=6773
  init_duration_map["a0b745bfcecce6c5dfa80d8e5302cc2e1ffbf9af"]=179
  init_duration_map["c5ae61d6d75b1948b71f56d49383c42ffe5ac684"]=402
  init_duration_map["c5e35a591b9915b413944fa520c63c7d62433d6a"]=115
  init_duration_map["d3c363b364961902abd9a0b78696201c8224b79e"]=18754
  init_duration_map["9d1f3cce4ca7d5801a3f961c4acfa440f89e20b1"]=1059
  init_duration_map["d3a3f36e6aa8caf20215fcb8846c0e3a28192dd2"]=192
  init_duration_map["695cf76d638a58ab7fad5f3789f923aa1ab5c08b"]=200
  init_duration_map["f24ce71ecf08e7025888b82344a30eb3659db9e4"]=238
  init_duration_map["1eeb44881951fddba002f6ceac52c324bacc46ce"]=170
  init_duration_map["3c66fdfa3354c35841d82814ec92a826a26f5e24"]=131
  init_duration_map["9c1222b7b54433655b7c623555490f47042cafdc"]=8219
  init_duration_map["96bfbddf7a5d83b8db82d3da94f59ce43448c70b"]=2294
  init_duration_map["167e8222a8cb9f19dbc77c280bf0f7fc0d8bb9b6"]=835
  init_duration_map["7c6e126925726e46adc7fd07d609291151fd010a"]=161
  init_duration_map["437c4d38a98713022a248c909a7c3f8a0a83df38"]=114
  init_duration_map["843fd86a7b09f0fe324f8c1a505e915eb8db7e0c"]=4430
  init_duration_map["1de3e81af891eaa431ef3ab4d857c5fe46140d6b"]=7909
  init_duration_map["9cdb9644ccd8e9a99f4d95c7a3e4bdebb181e717"]=750
  init_duration_map["39476527c2c0688970b71bffd453efd386a5e7f1"]=6505
  init_duration_map["ba815b764ce640a42ed3fc8e62108086a3527749"]=4094
  init_duration_map["12475a9fcc66e46025459323022765512f36bb5e"]=1401
  init_duration_map["e48366b1aa5c42a01ff12a6f82df8d847811eb02"]=8902
  init_duration_map["336ed51d1ab1cd7dd184d3362fe62b075b75693f"]=8056
  init_duration_map["7758df61747db33f0af9d8f33d3c985c5465a8be"]=892
  init_duration_map["a971cbd8b5ef198a52ba9fea819f5a2166c07131"]=6885
  init_duration_map["02a1fde9f06dbc17079fd3f45fca8d5ed03b8e35"]=187
  init_duration_map["6f83e25d1fad0b50fe5434423db84731f9a166c2"]=2524
  init_duration_map["4505ee1ca189e2eb39f5992d5834440c90a2809b"]=268
  init_duration_map["0a0ca59b77959302611c57d75ffc599485c12e34"]=7860
  init_duration_map["e9f83d135d0e7a1b8c72dd3e9adf54e0b4f99b29"]=790
  init_duration_map["b9eee161b9355279824f7ff806f8b3ebd90adeaa"]=135
  init_duration_map["a0c7da07974f39ffb7c3aecf29e41bbadb61942f"]=3292
  init_duration_map["030707aacf589606b1160912bdf4396b2d252915"]=811
  init_duration_map["a7475c9189f0075180926fdad4fe515e90dfeca6"]=1087
  init_duration_map["432db71e1c7a682dd3b9e246650eb669bc47eca9"]=212
  init_duration_map["037ee604ba07d956d1e17f65f74ce2c8a104be66"]=150
  init_duration_map["076d225b93ba9442adc0c88325390ca5868ed37e"]=1413
  init_duration_map["62aeb199e061f786cf94541daddaca9c225db708"]=198
  init_duration_map["0e34c8dd3e11af8e111714a5e604f32af2aa4b75"]=1508
  init_duration_map["ce9a3976d666f967e1ae802f74ed51d58a2fee97"]=5967
  init_duration_map["bfe1ad9927f9fbcc1b9989d7fc82ca92fabfc6a4"]=158
  init_duration_map["c4ad822e4ee94b6d8e9daae0ac5899b2b3fd5fbc"]=343
  init_duration_map["241e17da836a748f390175d3a119a69c859660f5"]=133
  init_duration_map["02e6d032c352a70e31e0a9faffebf8b81e363234"]=2777
  init_duration_map["a72d5598180b6661659df944d953c1b10e43ac9e"]=723
  init_duration_map["478a845b90e2f829ea7ebe073faf257c6ad4e476"]=946
  init_duration_map["91dc0d70cdaf2c8776d67070a5ed1b4fc56d94cc"]=3867
  init_duration_map["7ba8ce84afa16f55aa4846ba831df1fe534eec39"]=9910
  init_duration_map["4f003a7743535e73a881723003499dd80ef5c3f3"]=1405
  init_duration_map["bcfda646c7117b789ca7a74877eef07e15941bf1"]=957
  init_duration_map["09940f8b95ae0351492eccab24247404ad9c1c69"]=4692
  init_duration_map["5b6d00d467537593c02dfe735e15c182056c572a"]=171
  init_duration_map["c98e37640ddfac6eee70ffa641bcb8356b4c2219"]=134
  init_duration_map["0181aff53a1108a9365bf812f9707a46a76a93ba"]=484
  init_duration_map["0ce048366c22cddd368f17658fc7a638fdb8dbc9"]=1099
  init_duration_map["6f8d9ac59bfa8eff8c53ddc83ad161b5372f14ef"]=138
  init_duration_map["0dd9b0ae5c58a00ac5bcc9b77100381444a854ad"]=340
  init_duration_map["f43c6f855fce936da9177b65dfc67ef579edf77e"]=2727
  init_duration_map["e9f5f7ff8346ada131bcf1074252bfe1c52c0e7e"]=207
  init_duration_map["85040824d05d037f172aaee016d37c83d502e084"]=114
  init_duration_map["4ea6b320928a7700a011553085d0011b03274bab"]=8116
  init_duration_map["e7379454c4e4fc248e57eebbbead54ad06d41cd2"]=9853
  init_duration_map["dd91924ca522df94bfbb19a135268e3ebf7b86d7"]=1951
  init_duration_map["02b0622f917d08394ca4deff84a5ad666e88a26b"]=190
  init_duration_map["ce9d9dee935af179e37766a9c4097f722312bbcf"]=4317
  init_duration_map["eb85e9b8f47351d5a4b5bd805a0c6331ce69f0d5"]=1926
  init_duration_map["3d08a8edbe8f51c67d2f6f5129098e9c4e019577"]=1566
  init_duration_map["a515e0378196b59412af9070160660351cad3553"]=7750
  init_duration_map["713401e8e946f27b03ec8dd51ddefab71207f20e"]=4418
  init_duration_map["8d2616d3037a928c1437408a5ca5c449a0881a70"]=137
  init_duration_map["ecededdec4d51fb740472ba33ba7fb13304edfb9"]=6963
  init_duration_map["2792e77fd40ef65be920b558674a890545f394f5"]=802
  init_duration_map["0b7aa77d990dd331e769a3129c26e4764b688681"]=1222
  init_duration_map["c48b5f84a3d3182a96ad72768ce8159b6ed39afb"]=344
  init_duration_map["259d76d439c1b54929d5eb64b9e64792ac780950"]=191
  init_duration_map["a22beea5881245fc842dbdec3b7fce85f8bdc12b"]=6100
  init_duration_map["6c354ee49b7a97c9c213ec55de3dfe4778cfd755"]=211
  init_duration_map["d1b069252d398defdadb07db78d73927e38715a8"]=151
  init_duration_map["8d601befc944a304728043992b5a0a638835956f"]=767
  init_duration_map["0eac8d22a27c6da323865039d81b53cbde0217dd"]=1718
  init_duration_map["cac3feff9bf6917beb7cde29602e56f6426feff7"]=1397
  init_duration_map["ac3f183748d38c4d741aee5bf1ca607479fcd557"]=138
  init_duration_map["280c68a3dbe9b7936e81d0927ded51437e1cfb20"]=1502
  init_duration_map["e19ebb5ea36e72a08c970486ea684f3afa3d72f7"]=664
  init_duration_map["663e9d923bb282b4624f45707582ff2e33666a1b"]=4148
  init_duration_map["13a319213ec8529389b4bb85c3e829d5f4f62ed8"]=1225
  init_duration_map["635388c9afcb28ada1cf9e3edfac8558ab68d7ea"]=9819
  init_duration_map["db697eb65ea626e8a02baeb261e315210891fea8"]=4209
  init_duration_map["92ea72217b73f1d711a7dd007def2a73abadecac"]=9640
  init_duration_map["a38f57ce86ee7712e5d3ea1eb62718b6be1f9c74"]=6171
  init_duration_map["2629d957eaee009d2da1336402bf3e43dab6d7ea"]=665
  init_duration_map["bda42b5d70c079e20da06c6ddf4f618dc1afe8ef"]=277
  init_duration_map["2d84caf37d0d3e437e97161ea835537c2b4e6c2e"]=210
  init_duration_map["32e2365ac01be7f166a2b9c8ee9d6c2093002e78"]=5642
  init_duration_map["79866642aa723fbee7417bd357d4b737f7983d50"]=2385
  init_duration_map["b7624d7917d5dbee72eb189137772afb4dbe1f85"]=1506
  init_duration_map["89d061dc51ccc16c0a8756084c61ea8a3e01eb0d"]=8460
  init_duration_map["34d0b8062b8b15e344b912ba8190ee9c4f5209b3"]=1547
  init_duration_map["3b688bdd12d280a9b7fd3518881a7251753a8b35"]=6154
  init_duration_map["9d7398776642f6deae913b889748a38577637077"]=5693
  init_duration_map["2a8a5d163c8167ae0239185de35d53964f7b62db"]=363
  init_duration_map["c5ec647cf1d1fe393d60d02ed01b702e1e4ae545"]=7118
  init_duration_map["2762696c02c752dec332c3ff517f159db5f835cc"]=738
  init_duration_map["e10a70c6996a8fc21e255d91ba695549b46ff02f"]=353
  init_duration_map["b335cf32e0346f806e506633cee03ba4094fbca6"]=3583
  init_duration_map["8d153d7f7c7b257039dcfbe71cc9f8fb5de08703"]=2778
  init_duration_map["620711a699ef91594017c11ba1c44ae9491a389c"]=868
  init_duration_map["2f2f4172d4a070b1fbf83d449df406e2c5e1cf89"]=3979
  init_duration_map["a313eb3ebcd78877988ed4c68cdfc986e78aa98f"]=180
  init_duration_map["1985e23a0856041b1173345c2a7229127f4025e2"]=8824
  init_duration_map["324657697c9ded3da4838dc5ab23c71f140f66a9"]=147
  init_duration_map["9543b9e40874abe031f13f750b0be62131e8bd88"]=103
  init_duration_map["a10cad3cb5ea178e84ef84f0eb53133a15fb022d"]=1311
  init_duration_map["16056628af6b93d5a292adf1b01e0d63b01094fc"]=1360
  init_duration_map["b0d6b73d327f18928d08db5176ccb50170fa771b"]=188
  init_duration_map["20952e985be6abf488a639a189491c8d4da36849"]=142
  init_duration_map["b21e5259c41a324c2010b1d187366d448e3584ca"]=393
  init_duration_map["2ca98c95d7fe3ff31beb3c77b0496c4962c18d15"]=1306
  init_duration_map["691e44923cf182f29b793be363ed27de32cdb274"]=2660
  init_duration_map["dc257b03db32f877738b4d6646a95f22ca588f13"]=4157
  init_duration_map["aa01f5cb0c72e80e14226d9351dfa6bd72e718ca"]=84
  init_duration_map["7b63bc5e5f64f462d4e0845713bf99a1154f9599"]=125
  init_duration_map["51df273902d3ce5c5bc36f216cdd2ac9bf66187f"]=208
  init_duration_map["78fb2956bdc87d380f7f3350520c2b190c52d026"]=1489
  init_duration_map["423eebeea3d2d8dd7d4c460c25420961a4338b02"]=7488
  init_duration_map["97256a15f348b0350a93129691ebdafd14227c9d"]=2961
  init_duration_map["737ca382d836d82c4e2dc2fc32f0f04677c62e67"]=1269
  init_duration_map["7d4c0227f884cb28bffac18fabaf6a588f7a6231"]=811
  init_duration_map["ef010ff0357d0bb78ba169ab520ee26fcc6c659d"]=2003
  init_duration_map["796886e0acc5e8e1632bbc5ce788819fd9bcb930"]=1262
  init_duration_map["f11d91a25894dbaf414c00568bbfa6aa4358d83b"]=4175
  init_duration_map["8039ed4fdfbfc943c089ecac39a2cdf5d4c14347"]=196
  init_duration_map["f2ce5b954bb55d1c30b47dc12c66cc565961fcd7"]=5142
  init_duration_map["289e147cdde87c1e1c578efff0323f3f52748f2b"]=199
  init_duration_map["6d4686926c7dd79b48b8586458ce4c4ec0d5e076"]=106
  init_duration_map["d5fbc8a9e6832108e3bd14504386638c63be80d9"]=148
  init_duration_map["18972d239f250062cdb2e64084105e72682efe3d"]=1310
  init_duration_map["c58c4f6ff45e2952e462afdfbe388f0bb8a4d665"]=683
  init_duration_map["2f2f2274fd205a673a83bd0ccecf00a08bf8712c"]=2482
  init_duration_map["35abc2d27b129b752c68e6e45c266194575eebc2"]=7130
  init_duration_map["5b578a6b1c1b9450d92041199609aa073868f308"]=6008
  init_duration_map["334fe1233fdba26648ff0ecac6b146b32e5f0891"]=634
  init_duration_map["b49748342822eb2944f0766dadaaf695dcb6f579"]=1506
  init_duration_map["2efb52e16877d2f3affbd3aecae667416fcade45"]=1393
  init_duration_map["3757ca10ef8744acbe4ea05fee5b64cbc0e33998"]=2041
  init_duration_map["f0ba6563cdd171aea76993edd4c42484ce788a16"]=458
  init_duration_map["1ce17f6886de49095c034380da8b0629c3942bb8"]=1493
  init_duration_map["24c58d8f3b4b85881067a16a02765ba5f7652a07"]=185
  init_duration_map["2a67e4bc44ee4f68824a33205388f5b84f3883f6"]=1320
  init_duration_map["d768204d1c2c73e20323a166da6dc8df95a13941"]=667
  init_duration_map["adc49a6fc949e5ad693785dd96c86b36f267688e"]=4341
  init_duration_map["d7d06a5e685e4c3465f036c8530e895030cc22e0"]=9408
  init_duration_map["2ecfc67d09e13ae95216378800ead137a419798f"]=225
}


func NewD3(metaData *model2.Meta, config *config.ConfigD3) Scaler {
	client, err := platform_client2.New(config.ClientAddr)
	if err != nil {
		log.Fatalf("client init with error: %s", err.Error())
	}
	scheduler := &SimpleD3{
		config:         config,
		metaData:       metaData,
		platformClient: client,
		mu:             sync.Mutex{},
		wg:             sync.WaitGroup{},
		instances:      make(map[string]*model2.Instance),
		idleInstance:   list.New(),
	}
	log.Printf("New scaler for app: %s is created", metaData.Key)
	scheduler.wg.Add(1)
	go func() {
		defer scheduler.wg.Done()
		scheduler.gcLoop()
		log.Printf("gc loop for app: %s is stoped", metaData.Key)
	}()

	return scheduler
}

func (s *SimpleD3) Assign(ctx context.Context, request *pb.AssignRequest, req_total int) (*pb.AssignReply, error) {
	start := time.Now()
	instanceId := uuid.New().String()
	defer func() {
		log.Printf("Assign, request id: %s, instance id: %s, cost %dms", request.RequestId, instanceId, time.Since(start).Milliseconds())
	}()
	log.Printf("Assign, request id: %s", request.RequestId)
	s.mu.Lock()
	if element := s.idleInstance.Front(); element != nil {
		instance := element.Value.(*model2.Instance)
		instance.Busy = true
		s.idleInstance.Remove(element)
		s.mu.Unlock()
		log.Printf("Assign, request id: %s, instance %s reused", request.RequestId, instance.Id)
		instanceId = instance.Id
		return &pb.AssignReply{
			Status: pb.Status_Ok,
			Assigment: &pb.Assignment{
				RequestId:  request.RequestId,
				MetaKey:    instance.Meta.Key,
				InstanceId: instance.Id,
			},
			ErrorMessage: nil,
		}, nil
	}
	s.mu.Unlock()

	//Create new Instance
	resourceConfig := model2.SlotResourceConfig{
		ResourceConfig: pb.ResourceConfig{
			MemoryInMegabytes: request.MetaData.MemoryInMb,
		},
	}
	slot, err := s.platformClient.CreateSlot(ctx, request.RequestId, &resourceConfig)
	if err != nil {
		errorMessage := fmt.Sprintf("create slot failed with: %s", err.Error())
		log.Printf(errorMessage)
		return nil, status.Errorf(codes.Internal, errorMessage)
	}

	meta := &model2.Meta{
		Meta: pb.Meta{
			Key:           request.MetaData.Key,
			Runtime:       request.MetaData.Runtime,
			TimeoutInSecs: request.MetaData.TimeoutInSecs,
		},
	}
	instance, err := s.platformClient.Init(ctx, request.RequestId, instanceId, slot, meta)
	if err != nil {
		errorMessage := fmt.Sprintf("create instance failed with: %s", err.Error())
		log.Printf(errorMessage)
		return nil, status.Errorf(codes.Internal, errorMessage)
	}

	//add new instance
	s.mu.Lock()
	instance.Busy = true
	s.instances[instance.Id] = instance
	s.mu.Unlock()
	log.Printf("request id: %s, instance %s for app %s is created, init latency: %dms", request.RequestId, instance.Id, instance.Meta.Key, instance.InitDurationInMs)
  s.assign_time = int64(0.2 * float64(s.assign_time) + 0.8 * float64(instance.InitDurationInMs))

	return &pb.AssignReply{
		Status: pb.Status_Ok,
		Assigment: &pb.Assignment{
			RequestId:  request.RequestId,
			MetaKey:    instance.Meta.Key,
			InstanceId: instance.Id,
		},
		ErrorMessage: nil,
	}, nil
}

func (s *SimpleD3) Idle(ctx context.Context, request *pb.IdleRequest) (*pb.IdleReply, error) {
	if request.Assigment == nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("assignment is nil"))
	}
	reply := &pb.IdleReply{
		Status:       pb.Status_Ok,
		ErrorMessage: nil,
	}
	start := time.Now()
	instanceId := request.Assigment.InstanceId
	defer func() {
		log.Printf("Idle, request id: %s, instance: %s, cost %dus", request.Assigment.RequestId, instanceId, time.Since(start).Microseconds())
	}()
	//log.Printf("Idle, request id: %s", request.Assigment.RequestId)
	needDestroy := false
	slotId := ""
	if request.Result != nil && request.Result.NeedDestroy != nil && *request.Result.NeedDestroy {
		needDestroy = true
	}
	defer func() {
		if needDestroy {
			s.deleteSlot(ctx, request.Assigment.RequestId, slotId, instanceId, request.Assigment.MetaKey, "bad instance")
		}
	}()
	log.Printf("Idle, request id: %s", request.Assigment.RequestId)
	s.mu.Lock()
	defer s.mu.Unlock()
	if instance := s.instances[instanceId]; instance != nil {
		slotId = instance.Slot.Id
		instance.LastIdleTime = time.Now()
		if needDestroy {
			log.Printf("request id %s, instance %s need be destroy", request.Assigment.RequestId, instanceId)
			return reply, nil
		}

		if instance.Busy == false {
			log.Printf("request id %s, instance %s already freed", request.Assigment.RequestId, instanceId)
			return reply, nil
		}
		instance.Busy = false
		s.idleInstance.PushFront(instance)
	} else {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("request id %s, instance %s not found", request.Assigment.RequestId, instanceId))
	}
	return &pb.IdleReply{
		Status:       pb.Status_Ok,
		ErrorMessage: nil,
	}, nil
}

func (s *SimpleD3) deleteSlot(ctx context.Context, requestId, slotId, instanceId, metaKey, reason string) {
	log.Printf("start delete Instance %s (Slot: %s) of app: %s", instanceId, slotId, metaKey)
	if err := s.platformClient.DestroySLot(ctx, requestId, slotId, reason); err != nil {
		log.Printf("delete Instance %s (Slot: %s) of app: %s failed with: %s", instanceId, slotId, metaKey, err.Error())
	}
}

func (s *SimpleD3) gcLoop() {
	log.Printf("gc loop for app: %s is started", s.metaData.Key)
	ticker := time.NewTicker(s.config.GcInterval)
	for range ticker.C {
		for {
			s.mu.Lock()
			if element := s.idleInstance.Back(); element != nil {
				instance := element.Value.(*model2.Instance)
				idleDuration := time.Now().Sub(instance.LastIdleTime)
        gcThread := s.config.IdleDurationBeforeGC
        _, ok := init_duration_map[s.metaData.Key]
        if(ok) {
          gcThread = 5 * time.Minute 
        }
				if idleDuration > gcThread {
					//need GC
					s.idleInstance.Remove(element)
					delete(s.instances, instance.Id)
					s.mu.Unlock()
					go func() {
						reason := fmt.Sprintf("Idle duration: %fs, excceed configured duration: %fs", idleDuration.Seconds(), s.config.IdleDurationBeforeGC.Seconds())
						ctx := context.Background()
						ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
						defer cancel()
						s.deleteSlot(ctx, uuid.NewString(), instance.Slot.Id, instance.Id, instance.Meta.Key, reason)
					}()

					continue
				}
			}
			s.mu.Unlock()
			break
		}
	}
}

func (s *SimpleD3) Stats() Stats {
	s.mu.Lock()
	defer s.mu.Unlock()
	return Stats{
		TotalInstance:     len(s.instances),
		TotalIdleInstance: s.idleInstance.Len(),
	}
}

