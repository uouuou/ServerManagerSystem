(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-8ba4c450"],{"19f3":function(t,n,e){},"977e":function(t,n,e){"use strict";e("19f3")},c27c:function(t,n,e){"use strict";e.r(n);var o=function(){var t=this,n=t.$createElement,e=t._self._c||n;return e("div",[e("el-row",{staticClass:"nat-content",attrs:{gutter:20}},[e("el-col",{attrs:{span:8}},[e("div",{staticClass:"nat-content-item"},[e("span",[t._v("FRP线上版本号：")]),t._v(" "+t._s(t.source.frp_online_version)+" ")])]),e("el-col",{attrs:{span:8}},[e("div",{staticClass:"nat-content-item"},[e("span",[t._v("是否启用FRP：")]),t._v(" "+t._s(1===t.source.frp_run?"启用":2===t.source.frp_run?"不启用":"")+" "),e("el-switch",{staticClass:"switch-set",attrs:{"active-value":1,"inactive-value":2},on:{change:t.switchConfig},model:{value:t.config,callback:function(n){t.config=n},expression:"config"}})],1)]),e("el-col",{attrs:{span:8}},[e("div",{staticClass:"nat-content-item"},[e("span",[t._v("服务端版本号：")]),t._v(" "+t._s(t.source.frp_version)+" ")])])],1),e("div",{staticClass:"nat-button"},[e("el-button",{attrs:{type:"primary",plain:""},on:{click:t.getInfo}},[t._v("线上FRP信息")]),e("el-button",{attrs:{type:"primary",plain:""},on:{click:t.clickUpdate}},[t._v("更新FRP")]),e("el-button",{attrs:{type:"primary",plain:""},on:{click:t.setFile}},[t._v("设置FRPS配置文件")])],1),t.codeData?e("table",{staticClass:"nat-table"},t._l(t.codeData,(function(n,o,a){return e("tr",{key:a},[e("td",[t._v(t._s(o))]),e("td",[t._v(t._s(n))])])})),0):t._e(),t.addFunDialog.visible?e("FrpForm"):t._e()],1)},a=[],r=e("5530"),c=e("2f62"),i=function(){var t=this,n=t.$createElement,e=t._self._c||n;return e("el-dialog",{attrs:{title:t.addFunDialog.title,visible:t.addFunDialog.visible,width:"40%"},on:{"update:visible":function(n){return t.$set(t.addFunDialog,"visible",n)},close:t.closeDialog}},[e("el-form",{ref:"form",attrs:{model:t.form,"label-width":"100px",autocomplete:"off"}},[e("el-form-item",{attrs:{label:"配置文件",prop:"data",rules:[{required:!0,message:"该项不能为空，请填写完整信息",trigger:"blur"}]}},[e("el-input",{attrs:{type:"textarea",rows:"6",placeholder:"请输入配置文件"},model:{value:t.form.data,callback:function(n){t.$set(t.form,"data",n)},expression:"form.data"}})],1)],1),e("div",{staticClass:"dialog-footer",attrs:{slot:"footer"},slot:"footer"},[e("el-button",{on:{click:t.closeDialog}},[t._v("取 消")]),e("el-button",{attrs:{type:"primary",loading:t.submitLoading},on:{click:function(n){return t.submit("form")}}},[t._v("确 定")])],1)],1)},s=[],u=e("c74b"),f={name:"dialog-form",data:function(){return{submitLoading:!1,form:{data:void 0}}},computed:Object(r["a"])({},Object(c["b"])(["addFunDialog"])),created:function(){var t=this;Object(u["d"])().then((function(n){t.form.data=n.data}),(function(n){t.$message.error(n.message)}))},methods:{closeDialog:function(){var t={row:"",type:"",title:"",visible:!0};this.$store.commit("SET_DIALOG",t)},submit:function(t){var n=this;this.$refs[t].validate((function(t){if(t){n.submitLoading=!0;var e=Object.assign({},n.form);console.log(e),Object(u["k"])(e).then((function(t){n.$message.success(t.message),n.closeDialog(),n.submitLoading=!1}),(function(t){n.$message.error(t.message),n.submitLoading=!1}))}}))}}},l=f,d=e("2877"),b=Object(d["a"])(l,i,s,!1,null,null,null),p=b.exports,m={name:"NPS",computed:Object(r["a"])({},Object(c["b"])(["addFunDialog"])),components:{FrpForm:p},data:function(){return{source:"",codeData:"",config:""}},created:function(){this.sourceLoad()},methods:{sourceLoad:function(){var t=this;Object(u["a"])().then((function(n){var e=n.data;t.config=e.frp_run,t.source=e}),(function(n){t.$message.error(n.message)}))},getInfo:function(){var t=this;this.$confirm("此操作将获取线上FRP信息, 是否继续?","提示",{type:"warning"}).then((function(){Object(u["b"])().then((function(n){var e=n.data;t.codeData=e}),(function(n){t.$message.error(n.message)}))})).catch((function(){}))},clickUpdate:function(){var t=this;this.$confirm("此操作将更新FRP, 是否继续?","提示",{type:"warning"}).then((function(){Object(u["c"])().then((function(n){var e=n.data;t.codeData=e}),(function(n){t.$message.error(n.message)}))})).catch((function(){}))},switchConfig:function(t){var n=this,e={frp_run:t,id:this.source.id};Object(u["i"])(e).then((function(t){n.$message.success(t.message),n.source.frp_run=t.data.frp_run,n.config=t.data.frp_run}),(function(t){n.$message.error(t.message)}))},setFile:function(){var t={row:{},type:void 0,title:"设置FRPS配置文件"};this.$store.commit("SET_DIALOG",t)}}},v=m,g=(e("977e"),Object(d["a"])(v,o,a,!1,null,null,null));n["default"]=g.exports},c74b:function(t,n,e){"use strict";e.d(n,"a",(function(){return a})),e.d(n,"b",(function(){return r})),e.d(n,"c",(function(){return c})),e.d(n,"f",(function(){return i})),e.d(n,"g",(function(){return s})),e.d(n,"h",(function(){return u})),e.d(n,"i",(function(){return f})),e.d(n,"d",(function(){return l})),e.d(n,"k",(function(){return d})),e.d(n,"j",(function(){return b})),e.d(n,"e",(function(){return p})),e.d(n,"l",(function(){return m}));var o=e("751a"),a=function(t){return Object(o["b"])("/v1/nat/get_conf",{fun:"frp"})},r=function(t){return Object(o["b"])("/v1/nat/frp/online")},c=function(t){return Object(o["b"])("/v1/nat/frp/update")},i=function(t){return Object(o["b"])("/v1/nat/get_conf",{fun:"nps"})},s=function(t){return Object(o["b"])("/v1/nat/nps/online",t)},u=function(t){return Object(o["b"])("/v1/nat/nps/update",t)},f=function(t){return Object(o["c"])("/v1/nat/set_conf?fun=frp",t)},l=function(t){return Object(o["b"])("/v1/nat/frp/get_conf",t)},d=function(t){return Object(o["c"])("/v1/nat/frp/set_conf",t)},b=function(t){return Object(o["c"])("/v1/nat/set_conf?fun=nps",t)},p=function(t){return Object(o["b"])("/v1/nat/nps/get_conf",t)},m=function(t){return Object(o["c"])("/v1/nat/nps/set_conf",t)}}}]);