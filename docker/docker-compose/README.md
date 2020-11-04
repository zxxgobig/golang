# 开发环境部署说明

## 工具
先要自行安装`docker` `docker-compose` `git`

## 推荐的本地代码目录结构
```
code
├── be
├── tob-ats
├── fe
├── f2e
│ ├── common
│ └── public
└── thirdsrc
└── conf
└── log
```
__别忘了log目录,代码中无法自建__

## 项目依赖
### 本仓库作为E成业务的重要组成部分，整体业务的运行还依赖以下几个后端项目
- **be** `git@gitlab.ifchange.com:tob/web/be.git`
### 以及以下几个前端项目
- **fe** `git@gitlab.ifchange.com:tob/web/fe.git`
- **common** `git@gitlab.ifchange.com:tob/f2e/common.git`
- **public** `git@gitlab.ifchange.com:tob/f2e/public.git`

## 第三方库依赖
- **thirdsrc** 下载地址：`http://192.168.1.150:8090/download/attachments/65243101/thirdsrc.zip?version=1&modificationDate=1556444048000&api=v2`

## 额外配置依赖项
- 下载地址：`http://192.168.1.150:8090/download/attachments/65243101/conf.zip?version=1&modificationDate=1556443988000&api=v2`


## 集成环境运行
```
docker-compose build --no-cache --pull --force-rm
docker-compose up -d
```

## windows用户请注意！！！
> 因为docker中项目是通过文件映射的，windows下会有一个提示让用户确认共享磁盘（关注下任务栏），此时需要输入账号密码，但我们目前办公环境是没有设置win密码的，所以需要去设置一个密码，设置后开启共享，docker服务即可运行起来

## be项目额外配置
- be项目下`cp config/log4php_.properties.sample config/log4php_.properties`
- 新增be项目开发环境配置文件：`config/development/common.config.php`
```
<?php
$config['sess_cookie_name']             = 'IFCHANGE_TOB';
$config['sess_expiration']              = 60 * 60 * 24 * 30;
$config['sess_expire_on_close'] = FALSE;
$config['sess_encrypt_cookie']  = FALSE;
$config['sess_use_database']    = FALSE;
$config['sess_table_name']              = 'ci_sessions';
$config['sess_match_ip']                = FALSE;
$config['sess_match_useragent'] = TRUE;
$config['sess_time_to_update']  = 300;

$config['sess_storage']   = array(
	'development' => array('host'=>'192.168.1.201', 'port'=>'11212')
);

$config['cookie_domain']    = ".ifchange.com";

//开发环境域名
$config['base_url'] =  "http://dev.tob.ifchange.com/";
$config['appbase_url'] =  "https://dev.tob.ifchange.com/";

//开发环境简历下载域名
$config['download_domain'] =  "http://download.dev.ifchange.com/";

//开发C端环境域名
$config['toc_base_url'] =  "http://www.dev.cheng95.com/";

//注册域名
$config['reg_domain'] = 'http://liugj.easyhunter.dev.com/';

//充值网址
$config['charge_domain'] = 'http://pay.dev.ifchange.com/';

//礼包码
$config['package_domain'] = 'http://package.dev.ifchange.com/';

//二维码
$config['qrcode_domain'] = 'http://mengshuai.zhu.easyhunter.dev.com/';

//消息系统地址
$config['saleDomain'] = 'http://yanghua.dev.ifchange.com/';

//用户行为分析系统域名
$config['ubaDomain'] = 'http://pv.testing2.ifchange.com/';

//消息中心配置
$config['messageCenterAppId'] = 5;
$config['messageCenterToken'] = 'c6ebfd2da9fc16a779b7ba4925f1a265';

//邮件退订域名
$config['unsubscribeDomain'] = 'http://message.dev.ifchange.com/';

//支付中心配置
$config['payCenterToken'] = 'c6ebfd2da9fc16a779b7ba4925f1a265';

//礼包中心配置
$config['giftCenterToken'] = 'c6ebfd2da9fc16a779b7ba4925f1a265';

//静态文件地址
$config['staticPath'] = array(
	'//img.dev.ifchange.com'
);

//用户上传图片服务器地址
$config['imgServer'] = 'https://uimg.dev.ifchange.com/';

// 运营域名
$config['operationDomain'] = '';

$config[ 'wzpDomain' ] = 'http://m.dev3.ifchange.com/';
$config[ 'tocDomain' ] = 'http://zpp.m.zms.dev.cheng95.com/';
$config[ 'apptocDomain' ] = 'https://zpp.m.zms.dev.cheng95.com/';

// common下views目录的路径
$config['twiggy_dir'] = '/opt/wwwroot/tob/web/common/src/php-views/';

//双11活动时间
$config['activity_free_day']['start'] = 1447171200;
$config['activity_free_day']['end'] = 1447862400;

$config['zookeeperAddress'] = '192.168.1.200:2181';
$config['zookeeperRoot'] = '/opt/wwwroot/rd/2b/IfChangeTOB/be/third_party/gen-php';
$config['defaultDataDealTime'] = 1;
//意向邀约URL地址
$config['purposeInviteEmailUrl'] =  "http://guohui.dev.ka.ifchange.com/";
$config['defaultPurposeInviteTime'] = 1;

// 约ta新ats对接
$config['yuetaGetPositionUrl'] = 'http://dev.tob.ifchange.com/atsng/atsInternal/getRelateIcdcPositionIds';

//不合适邮件模板
$config['unsubscribe_url'] = 'http://dev.m.ifchange.com/';
//C端不合适列表
$config['tocApiDomain'] = 'http://api.chenwei.dev.cheng95.com/';
$config['c_more_chance_url'] = 'http://www.dev.cheng95.com/';

//ka图片配置地址
$config['kaArchivesIcon'] = array();
$config['kaArchivesIcon'][81] = '/tob/img/ka/qq.png';
$config['defaultResumeDealTime'] = 1;
//简历最近有更新
$config['resumeCurrentUpdateTime'] = 1;

// 伯乐系统接口
$config['boleApiDomain'] = 'http://api.bole.zms.dev3.ifchange.com';

// 伯乐系统域名
$config['boleDomain'] = 'http://bole.zms.dev3.ifchange.com';

// 大易平台接口wsdl服务
$config['dayeeApiDomain'] = 'http://win8.wintalent.cn/wt/xwebservices/eChengResumeNumService?wsdl';

// 大易平台接口wsdl服务固定请求参数
$config['dayeeParameter'] = [
	'corpCode' => 'ehuatai',
	'userName' => 'ECXZ',
	'password' => 'SWGsiiSpynHDWrOK',
];

// toc伯乐邀请注册
$config['invitationBoleDomain'] = 'http://bole.testing2.ifchange.com/welcome/invitation?tob_uid=';

// 人脉内推toc域名
$config['interpolateContactsBoleDomain'] = 'http://bole.testing2.ifchange.com';
$config['interpolateContactsTocDomain'] = 'http://www.testing2.cheng95.com';

// 伯乐手机端域名
$config['mBoleDomain'] = 'http://m.bole.testing2.ifchange.com';

// 查看人脉雷达限定帐号uid
$config['viewInterpolateRadarUid'] = 81;

// 配置微众银行主帐号id
$config['webankParentId'] = [];

// 配置万科主帐号id
$config['vankeParentId'] = [81];

// 配置联影帐号列表
$config['unitedImagingAccount'] = [100];

// 配置vanke接口请求域名
// $config['vankeDomain'] = 'http://120.77.168.249:8080/ehr-dic-service';
// $config['vankeDomain'] = 'http://120.77.168.249/ehr-dic-service-pre';
$config['vankeDomain'] = 'http://lms-app.vanke.com/ehr-dic-service-pre';

// vanke官网域名
$config[ 'vankeTocDomain' ] = 'http://m.vanke.testing2.cheng95.com';

// 免限制白名单
$config['interpolateContactsMoneyAccountLimit'] = [125981, 121553, 127996, 129057, 123982];// 太平人寿保险有限公司上海分公司, 上海众正医药科技有限公司, 深圳市雅诺信珠宝首饰有限公司, 深圳市学而思培训中心, 上海必胜客有限公司

// 测评配置
$config['estimate'] = [
	'nmsd' => [ // 诺姆四达
		'key' => 'n16731510709815296',
		'getActivitysUrl' => 'http://gl.normstar.net:8082/ns-napm-web/joinNapm/getActivitys.do', // 查询场次
		'getAccountUrl' => 'http://gl.normstar.net:8082/ns-napm-web/joinNapm/getAccountTestUrl.do', // 创建考生
		'estimateCategoryTypeMap' => [ // 测评种类类型 10(11认知力 12领导力素质 13销售人员素质 14产品经理素质 15专业人员适职 16潜在管理者)
			[
				'estimate_category_type' => 11,
				'event' => 'n17185123766634496', // 测评场次
				'name' => '认知力测试',
				'introduction' => '适用于全行业的各种岗位，从类比推理能力、图形推理能力、逻辑推理能力、数字推理能力四个方面对员工智商进行评估。',
				'apply' => ['行业不限', '岗位不限', '人才不限'],
				'time_limit' => '32分',
				'topic_number' => '32',
				'estimate_cost' => 1,
			],
			[
				'estimate_category_type' => 12,
				'event' => 'n17195855572109312', // 测评场次
				'name' => '领导力素质测试',
				'introduction' => '适用于全行业的中层管理岗位，从能力、个性、动力三个维度全面立体的评价其在如经营管理、任务管理、团队管理等多个方面所需的通用素质。',
				'apply' => ['行业不限', '岗位不限', '管理层'],
				'time_limit' => '1时15分',
				'topic_number' => '175',
				'estimate_cost' => 3,
			],
			[
				'estimate_category_type' => 13,
				'event' => 'n17185152712312833', // 测评场次
				'name' => '销售人员素质测试',
				'introduction' => '适用于全行业的销售岗位，全面立体的评价销售人员在如产品销售、货款回收、客户关系建立等多个方面所需要的素质。',
				'apply' => ['行业不限', '岗位不限', '人才不限'],
				'time_limit' => '50分',
				'topic_number' => '175',
				'estimate_cost' => 2,
			],
			[
				'estimate_category_type' => 14,
				'event' => 'n17185117458794496', // 测评场次
				'name' => '产品经理素质测试',
				'introduction' => '适用于全行业产品管理岗，从能力、个性、动力三个层面体现了产品经理在产品需求分析、产品设计、分析市场趋势、团队管理等方面所需要的素质。',
				'apply' => ['行业不限', '产品类', '管理层'],
				'time_limit' => '2时29分',
				'topic_number' => '263',
				'estimate_cost' => 5,
			],
			[
				'estimate_category_type' => 15,
				'event' => 'n17195883069966336', // 测评场次
				'name' => '专业人员适职测试',
				'introduction' => '适用于全行业各种岗位，全面立体的评价高端专业技术人才诸如在国际视野、技术实现、突破创新等多个方面所需要的素质。',
				'apply' => ['行业不限', '岗位不限', '人才不限'],
				'time_limit' => '1时30分',
				'topic_number' => '188',
				'estimate_cost' => 8,
			],
			[
				'estimate_category_type' => 16,
				'event' => 'n17195875338946560', // 测评场次
				'name' => '潜在管理者测试',
				'introduction' => '适用于全行业所有管培生、管理层后补等人员，针对如思维策略、视野角度、领导个性和领导者动力等具体指标出发，去考察特定人才在管理方面所具备的潜在素质与能力。',
				'apply' => ['行业不限', '岗位不限', '人才不限'],
				'time_limit' => '1时10分',
				'topic_number' => '230',
				'estimate_cost' => 5,
			],
		],
	],
	'saville' => [ // Saville
		'key' => '1a0c4e9356ea43b6bfddacf4f5d7',
		'session' => 'https://eztest.org/tenant/api/session/', // 查询场次
		'entry' => 'https://eztest.org/tenant/api/session/%s/entry/', // 创建考生
		'end' => 'https://eztest.org/tenant/api/session/%s/entry/end/exam/', // 结束考生
		'createSession' => 'https://eztest.org/tenant/api/session/', // 创建场次
		'updateSession' => 'https://eztest.org/tenant/api/session/%s/', // 修改场次
		'form' => 'https://eztest.org/tenant/api/form/list/', // 查询试卷
		'estimateCategoryTypeMap' => [ // 测评种类类型 20(21职业倾向 22全方位个人 23情绪倾向 24面试指导)
			[
				'estimate_category_type' => 21,
				'event' => 13935, // 测评场次
				'name' => '职业倾向测试', // 专家
				'introduction' => '适用于全行业的所有岗位在职业发展中的求职规划、工作绩效预测、工作适应分析，工作领域相关的动机、偏好、需求和才能等相关信息预测。',
				'apply' => ['行业不限', '岗位不限', '人才不限'],
				'time_limit' => '25分',
				'topic_number' => '72',
				'estimate_cost' => 15,
			],
			[
				'estimate_category_type' => 22,
				'event' => 13935, // 测评场次
				'name' => '全方位个人测试',
				'introduction' => '帮助企业全面了解测试者的个性心理特征，如其性格、才能、动机的得分，胜任力潜能分析，文化/环境匹配预测等。',
				'apply' => ['行业不限', '岗位不限', '人才不限'],
				'time_limit' => '25分',
				'topic_number' => '72',
				'estimate_cost' => 10,
			],
			[
				'estimate_category_type' => 23,
				'event' => 13935, // 测评场次
				'name' => '情绪倾向测试',
				'introduction' => '适用于全行业的所有岗位的职场选拔环节，不仅能测量当前心理健康状况，还能预测未来心理健康的稳定程度。',
				'apply' => ['行业不限', '岗位不限', '人才不限'],
				'time_limit' => '25分',
				'topic_number' => '72',
				'estimate_cost' => 10,
			],
			[
				'estimate_category_type' => 24,
				'event' => 13935, // 测评场次
				'name' => '面试指导', // 胜任力
				'introduction' => '帮助企业显著减少面试准备时间，为直线经理和招聘者提供标准化问题，预测工作绩效和潜能，关注最能显著预测工作绩效的胜任力特征的等。',
				'apply' => ['行业不限', '岗位不限', '人才不限'],
				'time_limit' => '25分',
				'topic_number' => '72',
				'estimate_cost' => 10,
			],
		],
	],
	'ddi' => [ // DDI
		'webservice' => 'https://stg.services.ddiworld.com/integration/SOAP/HRXML24/0324196F-6654-4425-96AD-B0202FCB585D', // 创建考生
		'ProjectKey' => '0324196F-6654-4425-96AD-B0202FCB585D',
		'UserName' => 'eCh3ng!',
		'Password' => '6M5ee48ng\C0',
		'clientKey' => '96703048-2512-49c0-94d0-fa121574f418',
		'estimateCategoryTypeMap' => [ // 测评种类类型 30(31销售人员素质测试 32专业人员适职测试 33领导力素质测试 34认知力测试)
			[
				'estimate_category_type' => 31,
				'event' => '56BC5B80-CCCF-4BF1-A6AD-5EC1B4AC9AB4', // 测评场次
				'name' => '销售人员素质测试',
				'introduction' => '适用于全行业的销售岗位的识别与甄选，特别是顾问型销售人员，也适用于校园招聘销售类岗位。',
				'apply' => ['行业不限', '岗位不限', '人才不限'],
				'time_limit' => '60分',
				'topic_number' => '100+',
				'estimate_cost' => 3,
			],
			[
				'estimate_category_type' => 32,
				'event' => '085DBD35-1ED8-441C-9D70-ECE232C055E9', // 测评场次
				'name' => '专业人员适职测试',
				'introduction' => '适用于全行业的专业岗位，工作经验0~3年，如：财务，法务，人力资源，工程技术，研发等。也适用于校园招聘。',
				'apply' => ['行业不限', '岗位不限', '人才不限'],
				'time_limit' => '60分',
				'topic_number' => '100+',
				'estimate_cost' => 3,
			],
			[
				'estimate_category_type' => 33,
				'event' => '12D9CF11-1728-4ECE-80ED-C6C38826636B', // 测评场次
				'name' => '领导力素质测试',
				'introduction' => '适用于全行业和全职能领域的中基层领导岗位的识别与甄选，如：主任、经理、总监等。',
				'apply' => ['行业不限', '岗位不限', '人才不限'],
				'time_limit' => '60分',
				'topic_number' => '100+',
				'estimate_cost' => 10,
			],
			[
				'estimate_category_type' => 34,
				'event' => 'FC830B94-FC1A-4949-9DC8-A61D341B0F34', // 测评场次
				'name' => '认知力测试',
				'introduction' => '适用于全行业的各种岗位，特别是需要高效解决问题和快速学习能力的岗位，可用于有初级经验人才的识别与甄选，也适用于校园招聘。',
				'apply' => ['行业不限', '岗位不限', '人才不限'],
				'time_limit' => '60分',
				'topic_number' => '25+',
				'estimate_cost' => 3,
			],
		],
	],
	'epi' => [ // EPI
		'key' => '1a0c4e9356ea43b6bfddacf4f5d7',
		'session' => 'https://eztest.org/tenant/api/session/', // 查询场次
		'entry' => 'https://eztest.org/tenant/api/session/%s/entry/', // 创建考生
		'end' => 'https://eztest.org/tenant/api/session/%s/entry/end/exam/', // 结束考生
		'createSession' => 'https://eztest.org/tenant/api/session/', // 创建场次
		'updateSession' => 'https://eztest.org/tenant/api/session/%s/', // 修改场次
		'form' => 'https://eztest.org/tenant/api/form/list/', // 查询试卷
		'estimateCategoryTypeMap' => [
			[
				'estimate_category_type' => 71,
				'event' => 14546, // 测评场次
				'name' => '认知力测试',
				'introduction' => '适用于全行业的各种岗位，也适用于校园招聘，从类比推理能力、图形推理能力、逻辑推理能力、数字推理能力四个方面对员工智商进行评估。',
				'apply' => ['行业不限', '基层员工', '校园招聘'],
				'time_limit' => '50分',
				'topic_number' => '41',
				'estimate_cost' => 1,
			]
		],
	],
];

//直辖市id在职位和简历中映射关系
$config['cityIdsRefer'] = array(
	"10"=>"105",
	"2"=>"33",
	"3"=>"34",
	"23"=>"267"
);
//工作年限
$config['work_time'] = array(
	0 => '',
	1 => '0_1',
	2 => '1_3',
	3 => '3_5',
	4 => '5_10',
	5 => '10_'
);
//年龄
$config['select_ages'] = array(
	0 => '',
	1 => '0_21',
	2 => '22_25',
	3 => '26_27',
	4 => '28_30',
	5 => '31_40',
	6 => '40_'
);
//年龄
$config['select_gender'] = array(
	0 => '',
	1 => 'M',
	2 => 'F'
);
//期望月薪分布
$config['hope_money'] = array(
	0 => '',
	1 => '0_5',
	2 => '5_8',
	3 => '8_12',
	4 => '12_15',
	5 => '15_20',
	6 => '20_30',
	7 => '30_50',
	8 => '50_'
);

//人才库文件夹下载资源失效时长（秒）
$config['talent_dir_download_expired'] = 30 * 86400;
//微信appid
$config['weixin_appid'] = 'wxe1dd082102e5a8d8';
$config['weixin_appsecret'] = 'c28b42785846c6a88b7c024b0a87fb5e';

// 云学堂相关
$config['yxtBuyRequestUrl'] = 'http://api.yunxuetang.com.cn/mallapi/v1/ifchange/events';
$config['yxtBuyRequestKey'] = '72ac3614bfa840f7b07c0cd9111d0b21';

//特殊uid，所有我们拥有名称的简历都显示名称给这个uid
$config['resume_special_uid'] = [81];

//销售智能H5页面搜索简历详情回去固定uid
$config['H5_uid'] = 81;

// 配置微众银行主帐号id
$config['webankParentId'] = [];// [81];

// 配置万科主帐号id
$config['vankeParentId'] = [81, 93005];

//officeweb365动态预览链接
$config['officeweb365'] = 'http://ow365.cn/?i=12687&del=1&furl=http://dev.tob.ifchange.com/api/resume/downloadoriginal?frame=1&token=';

//微众账号导出报表API地址
$config['webankReportExportHost'] = 'ats_report.dev.ifchange.com';

//新联康定制简历导出excel帐号uid
$config['XlkUid'] = 81;

//hr微课问题建议邮件发送目标地址
$config['hr_lesson_email'] = 'haoming.chi@ifchange.com';

//离职员工库：能关注的离职员工总数
$config['preemployeeFocusNumber'] = 1000;

//8月ATS运营活动时间节点
$config['integral_1'] = 1500912000;//07-25 1502640000 08-14
$config['integral_2'] = 1503244800;//08-21
$config['integral_3'] = 1503849600;//08-28
$config['integral_4'] = 1504454400;//09-04
$config['integral_5'] = 1504972800;//09-11

//三级账号共用人才库的账号
$config['onlyOneTalentAccount'] = [];

$config['exiaobao_report'] = [
    'to_mails'   => 'zhengfei.xie@ifchange.com',
    'from_email' => 'zhengfei.xie@ifchange.com',
    'from_name'  => '谢正飞',
    'password'   => 'Xie2017',
    'report_to'  => 'zhengfei.xie@ifchange.com',
];

//tob简历服务host
$config['tobResumeService'] = 'http://127.0.0.1:8080/rest';

// 配置天珑主帐号id
$config['tinnoParentId'] = [81];

// 配置金元证券主帐号id
$config['goldstateParentId'] = [81];

// 配置帮助中心用户提问发送邮箱
$config['helper_center_email'] = 'congyi.shen@ifchange.com';

// 重构以后的ATS主投接口
$config['atsngUrl'] = 'http://dev.tob.ifchange.com/atsng/delivery/addDeliveryForBe';

// 重构以后的外网同步职位接口在testing3
$config['atsngPositionUrl'] = 'http://dev.tob.ifchange.com/atsng';

$config['newTalentTopId'] = [81];
$config['enVersionTopId'] = [81];

// google analytics ID
$config['gaid'] = 'UA-118331797-1';

$config['rpc'] = array(
	'tob'=>'http://dev.tob.rpc/',
	'bi'=>'http://dev.bi.rpc/',
	'toh'=>'http://dev.toh.rpc/',
	'icdc'=>'http://dev.icdc.rpc/',
	'grab'=>'http://dev.grab.rpc/',
	'gsystem'=>'http://dev.gsystem.rpc/',
	'grabmail'=>'http://dev.grabmail.rpc/',
	'position'=>'http://dev.position.rpc/',
	'positiontoh'=>'http://dev.positiontoh.rpc/',
    'algorithm'=>'http://dev.algo.rpc/',
	'parser'=>'http://dev.nlpparser.rpc/tob',
    'parser_off'=>'http://dev.nlpparser.rpc/tob',
    'databustob'=>'http://192.168.1.200:51701/data-app-tob/tob-dispatch/dispatcher',
    'common-ats'=>'http://dev.tob.ifchange.com/',
);

// 产业BI域名
$config['industryInsight'] = '//insight.testing2.ifchange.com';
```

## Hosts
```
# ---------------- ifchange ----------------------
## swagger
192.168.3.144 dev.position.com

## 业务架构后台
192.168.2.244 dev.tob-arch-admin.ifchange.com

##  fe
#192.168.2.66 tob.dz.ifchange.com
#192.168.2.66 img.xt.ifchange.com

## web
127.0.0.1 tob.dz.ifchange.com
192.168.2.66 img.lg.ifchange.com
192.168.2.66 img.xt.ifchange.com

192.168.1.108 dev.position.rpc
192.168.1.108 dev.icdc.rpc
192.168.1.109 dev.toh.rpc
192.168.1.199 dev.tob.rpc

## 环境搭建引导
192.168.2.66 dev3.tob.ifchange.com
192.168.2.66 img.dev.ifchange.com
192.168.2.66 img1.dev.ifchange.com
192.168.2.66 img2.dev.ifchange.com
192.168.2.66 img3.dev.ifchange.com
192.168.2.66 img4.dev.ifchange.com

192.168.2.66 debug.testing3.ifchange.com
192.168.2.66 debug1.testing3.ifchange.com
192.168.2.66 debug2.testing3.ifchange.com
192.168.2.66 debug3.testing3.ifchange.com
192.168.2.66 debug4.testing3.ifchange.com

#192.168.2.66 img.testing3.ifchange.com
#192.168.2.66 img1.testing3.ifchange.com
#192.168.2.66 img2.testing3.ifchange.com
#192.168.2.66 img3.testing3.ifchange.com
#192.168.2.66 img4.testing3.ifchange.com

192.168.1.201 static.feng.liu.dev.com
127.0.0.1 dev.tob.ifchange.com
127.0.0.1 dz.tob.ifchange.com
127.0.0.1 cm.tob.ifchange.com
127.0.0.1 gaoyaoting.tob.ifchange.com
127.0.1.1 g.tob.ifchange.com
192.168.2.66 customize.tob.ifchange.com
#192.168.2.66 dev.tob.ifchange.com
192.168.1.236 zgc.tob.ifchange.com
127.0.0.1 yysdev.tob.ifchange.com
#127.0.0.1 img.dev.tob.ifchange.com
192.168.2.66 img.dev.tob.ifchange.com
127.0.0.1 img1.dev.tob.ifchange.com
127.0.0.1 img2.dev.tob.ifchange.com
127.0.0.1 img3.dev.tob.ifchange.com
192.168.1.110 api.bole.zms.dev3.ifchange.com
192.168.1.109 uimg.dev.ifchange.com
127.0.0.1 dev.yusheng.com
127.0.0.1 fastcgi.yusheng.com
127.0.0.1 dev.webank.ifchange.com
192.168.1.109 package.dev.ifchange.com
127.0.0.1 passport.ifchange.com
127.0.0.1 dev.external.ifchange.com
192.168.1.199 wenbin.ifchange.com
192.168.1.199 img.wenbin.ifchange.com
192.168.1.199 img1.wenbin.ifchange.com
192.168.1.199 img2.wenbin.ifchange.com
192.168.1.199 img3.wenbin.ifchange.com
192.168.1.199 www.dev.ifchange.com
192.168.1.199 img1.dev.ifchange.com
192.168.1.199 img2.dev.ifchange.com
192.168.1.199 img3.dev.ifchange.com
192.168.1.110 jw.www.zms.dev3.ifchange.com
192.168.1.199 img3.dev.ifchange.com
192.168.1.199 img3.dev.ifchange.com
192.168.1.199 img3.dev.ifchange.com
192.168.1.199 img3.dev.ifchange.com
192.168.1.201 yanghua.dev.ifchange.com

192.168.1.199 repos.ifchange.com
192.168.1.199 img.If.ifchange.com
192.168.1.199 If.tob.com


192.168.1.199 img.lf.ifchange.com

192.168.1.199 banner.dev.ifchange.com
192.168.1.199 2badmin.dev.ifchange.com

127.0.0.1 dev.position.com
127.0.0.1 dev.admin.customize.com
127.0.0.1 ats_report.dev.ifchange.com
192.168.2.66 _t_.tob.ifchange.com


#192.168.1.126 dev.ifchange.com
#192.168.1.126 img.dev.ifchange.com

192.168.1.110 txt.o-net.custom.yys.dev.cheng95.com
192.168.1.110 txt.o-net.custom.zl.dev.cheng95.com
192.168.1.110 g.jiangwei.dev.cheng95.com
192.168.1.110 jw.admin.zl.dev.cheng95.com
192.168.1.110 g.tangxuantao.dev.cheng95.com
192.168.1.110 g.wtt.dev.cheng95.com
192.168.1.110 uimg.dev.cheng95.com
192.168.1.110 txt.m.o-net.custom.zl.dev.cheng95.com
192.168.1.110 txt.m.o-net.custom.yys.dev.cheng95.com
192.168.1.110 txt.vanke.zl.dev.cheng95.com
192.168.1.110 txt.vanke.yys.dev.cheng95.com
192.168.1.110 txt.chinagreentown.yys.dev.cheng95.com
192.168.1.110 txt.salesdemo.yys.dev.cheng95.com
192.168.1.110 txt.ka-demo.yys.dev.cheng95.com
192.168.1.110 txt.m.ka-demo.yys.dev.cheng95.com
192.168.1.110 txt.m.salesdemo.yys.dev.cheng95.com
192.168.1.110 txt.m.icampus.yys.dev.cheng95.com
192.168.1.110 txt.m.iflytek.yys.dev.cheng95.com
192.168.1.110 txt.salesdemo.yangyusheng.dev.cheng95.com
192.168.1.110 txt.icampus.yangyusheng.dev.cheng95.com
192.168.1.110 txt.ka-demo.yangyusheng.dev.cheng95.com
192.168.1.110 txt.icampus.yys.dev.cheng95.com
192.168.1.110 txt.iflytek.yys.dev.cheng95.com
192.168.1.110 txt.jinke.yys.dev.cheng95.com
192.168.1.110 tangxuantao.icampus.yys.dev.cheng95.com
192.168.1.110 tangxuantao.salesdemo.yangyusheng.dev.cheng95.com
192.168.1.110 wtt.iflytek.zl.dev.cheng95.com
192.168.1.110 img.wutingting.dev.cheng95.com
192.168.1.110 g.wutingting.dev.cheng95.com
192.168.1.110 g1.wutingting.dev.cheng95.com

192.168.1.108 dev.position.rpc
192.168.1.108 dev.icdc.rpc
192.168.1.109 dev.toh.rpc
192.168.1.199 dev.tob.rpc
192.168.1.109 toh.kai.rpc
192.168.1.110 dev.toc.rpc
192.168.1.66 dev.gsystem.rpc

127.0.1.1 tob.dz.ifchange.com
127.0.1.1 lo.dz.ifchange.com
127.0.1.1 lo.tob.ifchange.com
127.0.0.1 partner-ka.dev.ifchange.com
192.168.1.139 img.wk.ifchange.com
192.168.2.66 img.lg.ifchange.com
192.168.2.66 img.zj.ifchange.com
192.168.2.133 dev.tob-arch-admin.ifchange.com
#192.168.2.66 partner.dz.ifchange.com
127.0.0.1 partner.dz.ifchange.com
192.168.2.66 img.xt.ifchange.com
192.168.1.109 pay.dev.ifchange.com

### 内网docker
docker.ifchange.com 211.148.28.11
# --------------- ifchange -----------------------
```

## 运行环境依赖(此节内容可不用关心)
### nginx
####nginx站点配置
- 见`./nginx/site.conf`

### PHP5.6.36
#### 注意以下扩展
- memcached
- msgpack
- gearman
- token_crypt 从`@张松`共享的目录`http://192.168.20.141/download/zips/php_token_crypt_nocert.zip`下载编译

#### 镜像构建文件
- 见`./php56/Dockerfile`

### PHP7.1.8
#### 注意以下扩展
- memcached
- msgpack
#### 镜像构建文件
- 见`./php/Dockerfile`


