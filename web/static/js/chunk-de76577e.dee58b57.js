(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-de76577e"],{"236b":function(e,t,n){"use strict";n("efee")},"65cb":function(e,t,n){"use strict";n.r(t);var r=function(){var e=this,t=e.$createElement,n=e._self._c||t;return n("div",[n("div",{staticClass:"hand-box"},[n("el-button",{attrs:{type:"primary"},on:{click:function(t){return e.editForm({},"")}}},[e._v("新增进程")]),n("el-dialog",{attrs:{title:"新增",visible:e.isVisible},on:{"update:visible":function(t){e.isVisible=t},close:e.closeDialog}},[n("el-form",{ref:"form",attrs:{model:e.form,"label-width":"90px"}},[n("input",{directives:[{name:"model",rawName:"v-model",value:e.form.id,expression:"form.id"}],attrs:{type:"hidden"},domProps:{value:e.form.id},on:{input:function(t){t.target.composing||e.$set(e.form,"id",t.target.value)}}}),n("el-form-item",{attrs:{prop:"name",label:"程序名:",rules:[{required:!0,message:"请填写内容",trigger:"blur"}]}},[n("el-input",{model:{value:e.form.name,callback:function(t){e.$set(e.form,"name",t)},expression:"form.name"}})],1),n("el-form-item",{attrs:{prop:"run_path",label:"路径:",rules:[{required:!0,message:"请填写内容",trigger:"blur"}]}},[n("el-input",{model:{value:e.form.run_path,callback:function(t){e.$set(e.form,"run_path",t)},expression:"form.run_path"}})],1),n("el-form-item",{attrs:{prop:"run_cmd",label:"运行命令:",rules:[{required:!0,message:"请填写内容",trigger:"blur"}]}},[n("el-input",{model:{value:e.form.run_cmd,callback:function(t){e.$set(e.form,"run_cmd",t)},expression:"form.run_cmd"}})],1),n("el-form-item",{attrs:{prop:"num",label:"启动次数:",rules:[{required:!0,message:"请填写内容",trigger:"blur"}]}},[n("el-input",{model:{value:e.form.num,callback:function(t){e.$set(e.form,"num",e._n(t))},expression:"form.num"}})],1),n("el-form-item",{attrs:{prop:"running",label:"是否启用:",rules:[{required:!0,message:"请选择",trigger:"change"}]}},[n("el-radio-group",{model:{value:e.form.running,callback:function(t){e.$set(e.form,"running",t)},expression:"form.running"}},[n("el-radio",{attrs:{label:1}},[e._v("启动")]),n("el-radio",{attrs:{label:2}},[e._v("关闭")])],1)],1),n("el-form-item",{attrs:{prop:"remark",label:"备注:",rules:[{required:!0,message:"请填写内容",trigger:"blur"}]}},[n("el-input",{attrs:{type:"textarea"},model:{value:e.form.remark,callback:function(t){e.$set(e.form,"remark",t)},expression:"form.remark"}})],1),n("el-form-item",[n("el-button",{attrs:{type:"primary"},on:{click:function(t){return e.submitForm("form")}}},[e._v("提交")])],1)],1)],1),n("el-dialog",{attrs:{title:"LOG查看",visible:e.isVisibleLog},on:{"update:visible":function(t){e.isVisibleLog=t},close:e.closeLog}},[n("div",{staticClass:"log-box"},e._l(e.contentLOG,(function(t,r){return n("div",{key:r,staticClass:"log-item"},[e._v(" "+e._s(t.msg)+" ")])})),0)])],1),n("el-table",{attrs:{data:e.sourceList,"row-key":"id","tree-props":{children:"children",hasChildren:"hasChildren"}}},[e._v("> "),n("el-table-column",{attrs:{prop:"id",label:"ID",width:"60"}}),n("el-table-column",{attrs:{prop:"name",label:"程序名"}}),n("el-table-column",{attrs:{prop:"run_path",label:"路径"}}),n("el-table-column",{attrs:{prop:"pid",label:"状态"}}),n("el-table-column",{attrs:{prop:"run_cmd",label:"运行命令","show-overflow-tooltip":""}}),n("el-table-column",{attrs:{prop:"remark",label:"备注","show-overflow-tooltip":""}}),n("el-table-column",{attrs:{prop:"update_at",label:"时间"}}),n("el-table-column",{attrs:{label:"操作",width:"300"},scopedSlots:e._u([{key:"default",fn:function(t){return[2===t.row.running?n("el-popconfirm",{attrs:{title:"确定要开启进程吗？"},on:{onConfirm:function(n){return e.processAction(t.row)}}},[n("el-button",{attrs:{slot:"reference",type:"warning",size:"mini",plain:""},slot:"reference"},[e._v("开启")])],1):n("el-popconfirm",{attrs:{title:"确定要开启关闭吗？"},on:{onConfirm:function(n){return e.processAction(t.row)}}},[n("el-button",{attrs:{slot:"reference",type:"warning",size:"mini",plain:""},slot:"reference"},[e._v("关闭")])],1),e._v(" "),1===t.row.running?[n("el-button",{attrs:{size:"mini"},on:{click:function(n){return e.logDataLoad(t.row)}}},[e._v("LOG")])]:e._e(),n("el-button",{attrs:{slot:"reference",type:"primary",plain:"",size:"mini"},on:{click:function(n){return e.editForm(t.row,"code_c_update")}},slot:"reference"},[e._v("修改")]),e._v(" "),n("el-popconfirm",{attrs:{title:"这是一段内容确定删除吗？"},on:{onConfirm:function(n){return e.itemDelete(t.row.id)}}},[n("el-button",{attrs:{slot:"reference",type:"danger",size:"mini",plain:""},slot:"reference"},[e._v("删除")])],1)]}}])})],1),n("pagination",{attrs:{total:e.total_amount,page:e.page,limit:e.page_size},on:{"update:page":function(t){e.page=t},"update:limit":function(t){e.page_size=t},pagination:e.handlePageChange}})],1)},o=[],s=n("fc23"),i=n("751a"),a=function(e){return Object(i["b"])("/v1/process/manage/process",e)},c=function(e){return Object(i["c"])("/v1/process/manage/process",e)},l=function(e){return Object(i["c"])("/v1/process/manage/process",e,"put")},u=function(e){return Object(i["c"])("/v1/process/manage/process",e,"delete")},p=function(e){return Object(i["c"])("/v1/process/manage/runProcess",e)},m=function(e){return Object(i["c"])("/v1/process/manage/offProcess",e)},f={name:"Dashboard",components:{Pagination:s["a"]},data:function(){return{sourceList:[],isVisible:!1,form:{},page_size:14,page:1,total_amount:0,type:"",isVisibleLog:!1,contentLOG:[],socket:"",queryItem:""}},created:function(){this.sourceLoad()},beforeDestroy:function(){this.close()},methods:{wsFn:function(){this.socket=new WebSocket("ws://"+window.location.host+"/api/open/ws_log"),this.socket.onopen=this.open,this.socket.onerror=this.error,this.socket.onmessage=this.wsMessage,this.socket.onsend=this.send},open:function(){this.send("连接成功","1")},error:function(){this.wsFn()},wsMessage:function(e){var t=JSON.parse(e.data);this.contentLOG.push(t)},close:function(){this.socket&&(this.socket.close(),this.socket="")},send:function(){var e={msg_type:1,token:sessionStorage.getItem("token"),log_file:this.queryItem.p_log};this.socket.send(JSON.stringify(e))},closeLog:function(){this.contentLOG=[],this.isVisibleLog=!1,this.socket.send(JSON.stringify({msg_type:2,token:sessionStorage.getItem("token")})),this.close()},logDataLoad:function(e){this.isVisibleLog=!0,this.queryItem=e,this.wsFn()},processAction:function(e){var t=this,n={id:e.id};1!==e.running?p(n).then((function(e){t.$message.success(e.message),t.sourceLoad()}),(function(e){t.$message.error(e.message)})):m(n).then((function(e){t.$message.success(e.message),t.sourceLoad()}),(function(e){t.$message.error(e.message)}))},handlePageChange:function(e){this.page=e.pageIndex,this.page_size=e.pageSize,this.sourceLoad()},itemDelete:function(e){var t=this;u({id:e}).then((function(e){t.$message.success(e.message),t.sourceLoad()}),(function(e){t.$message.error(e.message)}))},closeDialog:function(){this.isVisible=!1,this.form={}},submitForm:function(e){var t=this;this.$refs[e].validate((function(e){e&&(t.form.id?l(t.form).then((function(e){t.$message.success(e.message),t.sourceLoad(),t.closeDialog()}),(function(e){t.$message.error(e.message)})):c(t.form).then((function(e){t.$message.success(e.message),t.sourceLoad(),t.isVisible=!1,t.form={}}),(function(e){t.$message.error(e.message)})))}))},editForm:function(e,t){if(this.isVisible=!0,this.type=t,"code_c_add"===t&&(this.form.menu_code=e.menu_code,this.form.parent_code=e.menu_code),"code_c_update"===t){var n=Object.assign({},e);delete n.children,this.form=Object.assign({},e)}},sourceLoad:function(){var e=this,t={page_size:this.page_size,page:this.page};a(t).then((function(t){e.sourceList=t.data,e.total_amount=t.pages.total_amount}),(function(t){e.$message.error(t.message)}))}}},g=f,d=(n("a4a6"),n("2877")),h=Object(d["a"])(g,r,o,!1,null,null,null);t["default"]=h.exports},7156:function(e,t,n){var r=n("861d"),o=n("d2bb");e.exports=function(e,t,n){var s,i;return o&&"function"==typeof(s=t.constructor)&&s!==n&&r(i=s.prototype)&&i!==n.prototype&&o(e,i),e}},a4a6:function(e,t,n){"use strict";n("a9a2")},a9a2:function(e,t,n){},a9e3:function(e,t,n){"use strict";var r=n("83ab"),o=n("da84"),s=n("94ca"),i=n("6eeb"),a=n("5135"),c=n("c6b6"),l=n("7156"),u=n("c04e"),p=n("d039"),m=n("7c73"),f=n("241c").f,g=n("06cf").f,d=n("9bf2").f,h=n("58a8").trim,b="Number",_=o[b],v=_.prototype,k=c(m(v))==b,y=function(e){var t,n,r,o,s,i,a,c,l=u(e,!1);if("string"==typeof l&&l.length>2)if(l=h(l),t=l.charCodeAt(0),43===t||45===t){if(n=l.charCodeAt(2),88===n||120===n)return NaN}else if(48===t){switch(l.charCodeAt(1)){case 66:case 98:r=2,o=49;break;case 79:case 111:r=8,o=55;break;default:return+l}for(s=l.slice(2),i=s.length,a=0;a<i;a++)if(c=s.charCodeAt(a),c<48||c>o)return NaN;return parseInt(s,r)}return+l};if(s(b,!_(" 0o1")||!_("0b1")||_("+0x1"))){for(var w,L=function(e){var t=arguments.length<1?0:e,n=this;return n instanceof L&&(k?p((function(){v.valueOf.call(n)})):c(n)!=b)?l(new _(y(t)),n,L):y(t)},I=r?f(_):"MAX_VALUE,MIN_VALUE,NaN,NEGATIVE_INFINITY,POSITIVE_INFINITY,EPSILON,isFinite,isInteger,isNaN,isSafeInteger,MAX_SAFE_INTEGER,MIN_SAFE_INTEGER,parseFloat,parseInt,isInteger".split(","),z=0;I.length>z;z++)a(_,w=I[z])&&!a(L,w)&&d(L,w,g(_,w));L.prototype=v,v.constructor=L,i(o,b,L)}},efee:function(e,t,n){},fc23:function(e,t,n){"use strict";var r=function(){var e=this,t=e.$createElement,n=e._self._c||t;return n("div",{staticClass:"page-box"},[n("el-pagination",{attrs:{background:e.background,"current-page":e.currentPage,"page-size":e.limit,layout:e.layout,"page-sizes":e.pageSizes,total:e.total},on:{"update:currentPage":function(t){e.currentPage=t},"update:current-page":function(t){e.currentPage=t},"update:pageSize":function(t){e.limit=t},"update:page-size":function(t){e.limit=t},"size-change":e.handleSizeChange,"current-change":e.handleCurrentChange}})],1)},o=[],s=(n("a9e3"),{props:{total:{required:!0,type:Number},page:{type:Number,default:1},limit:{type:Number,default:20},pageSizes:{type:Array,default:function(){return[this.limit,20,30,50]}},layout:{type:String,default:"total, sizes, prev, pager, next, jumper"},background:{type:Boolean,default:!0},autoScroll:{type:Boolean,default:!0},hidden:{type:Boolean,default:!1}},computed:{currentPage:{get:function(){return this.page},set:function(e){this.$emit("update:page",e)}},pageSize:{get:function(){return this.limit},set:function(e){this.$emit("update:limit",e)}}},methods:{handleSizeChange:function(e){this.$emit("pagination",{pageIndex:1,pageSize:e})},handleCurrentChange:function(e){this.$emit("pagination",{pageIndex:e,pageSize:this.pageSize})}}}),i=s,a=(n("236b"),n("2877")),c=Object(a["a"])(i,r,o,!1,null,null,null);t["a"]=c.exports}}]);