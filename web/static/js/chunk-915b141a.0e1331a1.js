(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-915b141a"],{"519e":function(e,t,o){"use strict";o("d897")},"8e25":function(e,t,o){"use strict";o.r(t);var n=function(){var e=this,t=e.$createElement;e._self._c;return e._m(0)},s=[function(){var e=this,t=e.$createElement,o=e._self._c||t;return o("div",{staticClass:"xterm-box"},[o("div",{attrs:{id:"xterm"}})])}],i=(o("d3b7"),o("25f0"),o("fcf3")),r=o("47d0"),c=(o("173c"),o("abb2"),{name:"xterm",data:function(){return{term:"",socket:"",server_id:this.$route.query.id.toString(),lockReconnect:!1,timeout:28e3,timeoutObj:null,serverTimeoutObj:null,timeoutnum:null}},mounted:function(){this.init()},methods:{init:function(){var e=window.location.host;this.socket=new WebSocket("ws://"+e+"/api/open/ws"),this.socket.onopen=this.open,this.socket.onerror=this.error,this.socket.onmessage=this.wsMessage,this.socket.onsend=this.send},open:function(){this.initTerm(),this.send("","1")},error:function(){this.reconnect()},wsMessage:function(e){this.term.write(e.data)},close:function(){this.socket&&this.socket.close(),this.term&&this.term.dispose(document.getElementById("xterm"))},send:function(e,t){var o={msg_type:1==t?1:3==t?3:2,token:sessionStorage.getItem("token"),server_id:this.server_id,command:e,rows:this.xbH,cols:this.xbW};this.socket.send(JSON.stringify(o))},initTerm:function(){var e=document.querySelector(".xterm-box");this.xbH=Math.floor((e.offsetHeight-60)/17),this.xbW=Math.floor(document.querySelector(".xterm-box").offsetWidth/7);var t=this.term=new i["Terminal"]({fontFamily:'monaco, Consolas, "Lucida Console", monospace',rendererType:"canvas",rows:this.xbH,cols:this.xbW,convertEol:!0,disableStdin:!1,cursorStyle:"underline",cursorBlink:!0,theme:{foreground:"#7e9192",background:"#002833",cursor:"help",lineHeight:16}});t.open(document.getElementById("xterm"));var o=new r["FitAddon"];t.loadAddon(o),o.fit(),t.focus();var n=this;t.onData((function(e){n.send(e)}))}},beforeDestroy:function(){this.close(),clearTimeout(this.timeoutObj),clearTimeout(this.serverTimeoutObj),clearTimeout(this.timeoutnum)}}),u=c,a=(o("519e"),o("2877")),m=Object(a["a"])(u,n,s,!1,null,null,null);t["default"]=m.exports},d897:function(e,t,o){}}]);