(window.webpackJsonp=window.webpackJsonp||[]).push([[0],{131:function(e,t,a){e.exports=a(294)},294:function(e,t,a){"use strict";a.r(t);var n=a(24),r=a(25),i=a(27),s=a(26),o=a(28),l=a(0),c=a.n(l),d=a(16),m=a.n(d),h=a(129),p=a(38),u=a(18),f=a(19),g=a.n(f),v=a(32),b=a.n(v),w=a(35),E=a.n(w),y=a(23),k=a.n(y),C=a(34),S=a.n(C),x=a(46),N=a.n(x),j=a(29),O=a.n(j),D=a(114),I=a(5),z=a.n(I),A=a(115),T=a.n(A),L=a(118),B=a.n(L),F=a(119),P=a.n(F),R=a(122),K=a.n(R),q=a(120),H=a.n(q),M=a(121),J=a.n(M),V=a(62),W=a.n(V),_=a(78),U=a.n(_),G=a(117),X=a.n(G),$={success:T.a,warning:X.a,error:B.a,info:P.a};var Q=Object(u.withStyles)(function(e){return{success:{backgroundColor:H.a[600]},error:{backgroundColor:e.palette.error.dark},info:{backgroundColor:e.palette.primary.dark},warning:{backgroundColor:J.a[700]},icon:{fontSize:20},iconVariant:{opacity:.9,marginRight:e.spacing.unit},message:{display:"flex",alignItems:"center"}}})(function(e){var t=e.classes,a=e.className,n=e.message,r=e.onClose,i=e.variant,s=Object(D.a)(e,["classes","className","message","onClose","variant"]),o=$[i];return c.a.createElement(U.a,Object.assign({className:z()(t[i],a),"aria-describedby":"client-snackbar",message:c.a.createElement("span",{id:"client-snackbar",className:t.message},c.a.createElement(o,{className:z()(t.icon,t.iconVariant)}),n),action:[c.a.createElement(W.a,{key:"close","aria-label":"Close",color:"inherit",className:t.close,onClick:r},c.a.createElement(K.a,{className:t.icon}))]},s))}),Y=a(33),Z=a.n(Y),ee=a(45),te=a.n(ee),ae=a(47),ne=a.n(ae),re=function(e){function t(e){var a;return Object(n.a)(this,t),(a=Object(i.a)(this,Object(s.a)(t).call(this,e))).handleInfoClose=function(e,t){if("clickaway"!==t){var n=a.state.info;n.open=!1,a.setState({info:n})}},a.handleKeyDown=function(e){13===e.keyCode&&(a.handleRegister(),e.preventDefault())},a.state={info:{open:!1,variant:"success",message:""},infoOpen:!0,email:"",emailErr:!1,disabled:!1},a}return Object(o.a)(t,e),Object(r.a)(t,[{key:"handleEmail",value:function(e){var t=e.target.value;""!==t?""!==t&&(/^[a-z0-9]+([._\\-]*[a-z0-9])*@([a-z0-9]+[-a-z0-9]*[a-z0-9]+.){1,63}[a-z0-9]+$/.test(t)?this.setState({email:t,emailErr:!1}):this.setState({email:t,emailErr:!0})):this.setState({email:"",emailErr:!1})}},{key:"handleRegister",value:function(){var e=this;this.state.email?this.state.emailErr?this.setState({info:{open:!0,variant:"error",message:"\u90ae\u7bb1\u683c\u5f0f\u9519\u8bef"}}):(fetch("http://www.koogo.net:8080/user/auth/register",{method:"POST",headers:{"Content-Type":"application/json"},body:JSON.stringify({email:this.state.email})}).then(function(e){return 200===e.status?e.json():{code:e.status,message:"internal error"}}).then(function(t){if(200===t.code)return e.setState({disabled:!1,info:{open:!0,variant:"success",message:"\u6ce8\u518c\u6210\u529f, \u9a8c\u8bc1\u94fe\u63a5\u5df2\u53d1\u9001\u81f3\u60a8\u7684\u90ae\u7bb1, \u8bf7\u572830\u5206\u949f\u4ee5\u5185\u524d\u5f80\u5b8c\u6210\u6ce8\u518c!"}}),void e.setState({email:""});300!==t.code?301!==t.code?302!==t.code?e.setState({disabled:!1,info:{open:!0,variant:"error",message:"\u7cfb\u7edf\u9519\u8bef: "+t.message}}):e.setState({passwordErr:!0,disabled:!1,info:{open:!0,variant:"error",message:"\u9a8c\u8bc1\u90ae\u4ef6\u53d1\u9001\u5931\u8d25, \u8bf7\u7a0d\u540e\u518d\u8bd5"}}):e.setState({emailErr:!0,disabled:!1,info:{open:!0,variant:"error",message:"\u8be5\u90ae\u7bb1\u5df2\u88ab\u6ce8\u518c!"}}):e.setState({emailErr:!0,disabled:!1,info:{open:!0,variant:"error",message:"\u90ae\u7bb1\u683c\u5f0f\u4e0d\u6b63\u786e!"}})}).catch(function(t){alert(t),e.setState({disabled:!1,info:{open:!0,variant:"error",message:"\u7f51\u7edc\u9519\u8bef: "+t}})}),this.setState({disabled:!0})):this.setState({emailErr:!0,info:{open:!0,variant:"error",message:"\u90ae\u7bb1\u4e0d\u80fd\u4e3a\u7a7a!"}})}},{key:"render",value:function(){var e=this,t=this.props.classes;return c.a.createElement("div",{className:t.container},c.a.createElement(b.a,null),c.a.createElement(te.a,{hidden:!this.state.disabled,className:t.progress}),c.a.createElement(Z.a,{anchorOrigin:{vertical:"top",horizontal:"center"},open:this.state.info.open,autoHideDuration:2e3,onClose:this.handleInfoClose},c.a.createElement(Q,{onClose:this.handleInfoClose,variant:this.state.info.variant,message:this.state.info.message})),c.a.createElement(g.a,{className:t.paper},c.a.createElement(k.a,{component:"h1",variant:"h6",color:"primary",className:t.typography},c.a.createElement(N.a,{className:t.cloud}),"\u4e91\u8bb0"),c.a.createElement(S.a,{variant:"middle",className:t.divider}),c.a.createElement("form",{className:t.form,onKeyDown:this.handleKeyDown},c.a.createElement(O.a,{className:t.textField,id:"login-email",variant:"filled",label:"\u90ae\u7bb1",autoComplete:"off",value:this.state.email,onChange:function(t){return e.handleEmail(t)},onBlur:function(t){return e.handleEmail(t)},error:this.state.emailErr,disabled:this.state.disabled,required:!0}),c.a.createElement(E.a,{variant:"outlined",color:"primary",className:t.button,onClick:function(){return e.handleRegister()},disabled:this.state.disabled},"\u6ce8\u518c"),c.a.createElement(ne.a,{component:"button",variant:"body2",onClick:function(){e.props.history.push("/login")}},"\u5df2\u6ce8\u518c, \u767b\u9646"))))}}]),t}(c.a.Component),ie=Object(u.withStyles)(function(e){return{container:{height:"100vh",background:"linear-gradient(to right, #348AC7, #7474BF);",overflow:"hidden"},progress:{position:"fixed",width:"100%",background:"transparent"},paper:{maxWidth:400,margin:"50px auto",padding:"20px 0 40px 0"},typography:{textAlign:"center",verticalAlign:"middle",fontSize:"18px",lineHeight:"40px"},cloud:{verticalAlign:"top",fontSize:"35px",marginRight:10},divider:{margin:"0 0 40px 0"},form:{display:"flex",flexDirection:"column",justifyContent:"center",alignItems:"center"},textField:{width:"70%"},button:{fontSize:"18px",height:45,width:"40%",marginTop:20,marginBottom:5}}})(re),se=window.localStorage,oe=function(e){function t(e){var a;return Object(n.a)(this,t),(a=Object(i.a)(this,Object(s.a)(t).call(this,e))).handleInfoClose=function(e,t){if("clickaway"!==t){var n=a.state.info;n.open=!1,a.setState({info:n})}},a.handleKeyDown=function(e){13===e.keyCode&&(a.handleLogin(),e.preventDefault())},a.state={info:{open:!1,variant:"success",message:""},infoOpen:!0,email:"",password:"",emailErr:!1,passwordErr:!1,disabled:!1},a}return Object(o.a)(t,e),Object(r.a)(t,[{key:"componentDidMount",value:function(){se.token&&this.props.history.push("/profile")}},{key:"handleEmail",value:function(e){this.setState({email:e.target.value,emailErr:!1})}},{key:"handlePassword",value:function(e){this.setState({password:e.target.value,passwordErr:!1})}},{key:"handleLogin",value:function(){var e=this;this.state.email?this.state.password?(fetch("http://www.koogo.net:8080/user/auth/login",{method:"POST",headers:{"Content-Type":"application/json"},body:JSON.stringify({email:this.state.email,password:this.state.password})}).then(function(e){if(200===e.status){var t=e.headers.get("X-Login-token");return se.token=t,e.json()}return{code:e.status,message:"internal error"}}).then(function(t){if(200!==t.code)300!==t.code?303!==t.code?304!==t.code?e.setState({disabled:!1,info:{open:!0,variant:"error",message:"\u672a\u77e5\u9519\u8bef!"}}):e.setState({passwordErr:!0,disabled:!1,info:{open:!0,variant:"error",message:"\u5bc6\u7801\u9519\u8bef!"}}):e.setState({emailErr:!0,disabled:!1,info:{open:!0,variant:"error",message:"\u7528\u6237\u4e0d\u5b58\u5728!"}}):e.setState({emailErr:!0,passwordErr:!0,disabled:!1,info:{open:!0,variant:"error",message:"\u90ae\u7bb1\u6216\u5bc6\u7801\u683c\u5f0f\u4e0d\u6b63\u786e!"}});else{e.setState({disabled:!1,info:{open:!0,variant:"success",message:"\u767b\u9646\u6210\u529f!"}});var a=e;setTimeout(function(){a.props.history.push("/profile")},1e3)}}).catch(function(t){alert(t),e.setState({disabled:!1,info:{open:!0,variant:"error",message:"\u7f51\u7edc\u9519\u8bef: "+t}})}),this.setState({disabled:!0})):this.setState({passwordErr:!0,info:{open:!0,variant:"error",message:"\u5bc6\u7801\u4e0d\u80fd\u4e3a\u7a7a!"}}):this.setState({emailErr:!0,info:{open:!0,variant:"error",message:"\u90ae\u7bb1\u4e0d\u80fd\u4e3a\u7a7a!"}})}},{key:"render",value:function(){var e=this,t=this.props.classes;return c.a.createElement("div",{className:t.container,onSubmit:function(e){alert(e)}},c.a.createElement(b.a,null),c.a.createElement(te.a,{hidden:!this.state.disabled,className:t.progress}),c.a.createElement(Z.a,{anchorOrigin:{vertical:"top",horizontal:"center"},open:this.state.info.open,autoHideDuration:1e3,onClose:this.handleInfoClose},c.a.createElement(Q,{onClose:this.handleInfoClose,variant:this.state.info.variant,message:this.state.info.message})),c.a.createElement(g.a,{className:t.paper},c.a.createElement(k.a,{component:"h1",variant:"h6",color:"primary",className:t.typography},c.a.createElement(N.a,{className:t.cloud}),"\u4e91\u8bb0"),c.a.createElement(S.a,{variant:"middle",className:t.divider}),c.a.createElement("form",{className:t.form,onKeyDown:this.handleKeyDown},c.a.createElement(O.a,{className:t.textField,id:"login-email",variant:"filled",label:"\u90ae\u7bb1",autoComplete:"off",value:this.state.email,onChange:function(t){return e.handleEmail(t)},onBlur:function(t){return e.handleEmail(t)},error:this.state.emailErr,disabled:this.state.disabled,required:!0}),c.a.createElement(O.a,{className:t.textField,id:"login-password",variant:"filled",label:"\u5bc6\u7801",autoComplete:"new-password",type:"password",value:this.state.password,onChange:function(t){return e.handlePassword(t)},onBlur:function(t){return e.handlePassword(t)},error:this.state.passwordErr,disabled:this.state.disabled,required:!0}),c.a.createElement(E.a,{variant:"outlined",color:"primary",className:t.button,onClick:function(){return e.handleLogin()},disabled:this.state.disabled},"\u767b\u9646"),c.a.createElement(ne.a,{component:"button",variant:"body2",onClick:function(){e.props.history.push("/register")}},"\u6ca1\u6709\u8d26\u53f7, \u6ce8\u518c"))))}}]),t}(c.a.Component),le=Object(u.withStyles)(function(e){return{container:{height:"100vh",background:"linear-gradient(to right, #348AC7, #7474BF);",overflow:"hidden"},progress:{position:"fixed",width:"100%",background:"transparent"},paper:{maxWidth:400,margin:"50px auto",padding:"20px 0 40px 0"},typography:{textAlign:"center",verticalAlign:"middle",fontSize:"18px",lineHeight:"40px"},cloud:{verticalAlign:"top",fontSize:"35px",marginRight:10},divider:{margin:"0 0 40px 0"},form:{display:"flex",flexDirection:"column",justifyContent:"center",alignItems:"center"},textField:{width:"70%"},button:{fontSize:"18px",height:45,width:"40%",marginTop:20,marginBottom:5}}})(oe),ce=a(123),de=a.n(ce),me=a(124),he=a.n(me),pe=a(80),ue=a.n(pe),fe=a(63),ge=a.n(fe),ve=a(127),be=a.n(ve),we=a(128),Ee=a.n(we),ye=a(39),ke=a.n(ye),Ce=a(64),Se=a.n(Ce),xe=a(125),Ne=a.n(xe),je=a(126),Oe=a.n(je),De=window.localStorage,Ie=function(e){function t(e){var a;return Object(n.a)(this,t),(a=Object(i.a)(this,Object(s.a)(t).call(this,e))).state={avatar:"",id:"",email:"",nickname:"",info:{open:!1,variant:"success",message:""},disabled:!1},a}return Object(o.a)(t,e),Object(r.a)(t,[{key:"redirectLogin",value:function(){De.token="",this.props.history.push("/login")}},{key:"handleLogout",value:function(){var e=this;fetch("http://www.koogo.net:8080/user/auth/logout",{method:"POST",mode:"cors",headers:{"x-login-token":De.token}}).then(function(e){var t=e.headers.get("x-login-token");return t&&(De.token=t),e.json()}).then(function(t){De.token="",e.setState({disabled:!1,info:{open:!0,variant:"success",message:"\u6ce8\u9500\u767b\u9646\u6210\u529f!"}});var a=e;setTimeout(function(){a.props.history.push("/login")},1e3)}).catch(function(t){De.token="",e.setState({disabled:!1,info:{open:!0,variant:"error",message:"\u6ce8\u9500\u767b\u9646\u5931\u8d25!"}})})}},{key:"componentDidMount",value:function(){var e=this;""===De.token?this.props.history.push("/login"):fetch("http://www.koogo.net:8080/user/info/",{headers:{"x-login-token":De.token},mode:"cors"}).then(function(e){return e.json()}).then(function(t){200===t.code?e.setState({email:t.data.email,id:t.data.user_id,nickname:t.data.nickname,avatar:t.data.avatar}):(alert(t.message),e.redirectLogin())}).catch(function(t){alert(t),e.redirectLogin()})}},{key:"render",value:function(){var e=this,t=this.props.classes;return c.a.createElement("div",{className:t.root},c.a.createElement(b.a,null),c.a.createElement(Z.a,{anchorOrigin:{vertical:"top",horizontal:"center"},open:this.state.info.open,autoHideDuration:1e3,onClose:this.handleInfoClose},c.a.createElement(Q,{onClose:this.handleInfoClose,variant:this.state.info.variant,message:this.state.info.message})),c.a.createElement(de.a,{position:"static"},c.a.createElement(he.a,{variant:"dense",className:t.toolbar},c.a.createElement(ue.a,{className:t.wrapMenu},c.a.createElement(W.a,{className:t.menuButton,color:"inherit","aria-label":"Menu"},c.a.createElement(Ne.a,null)),c.a.createElement(k.a,{variant:"h6",color:"inherit"},"\u7528\u6237\u4fe1\u606f")),c.a.createElement(E.a,{color:"inherit",onClick:function(){e.handleLogout()}},c.a.createElement(Oe.a,null),"\u9000\u51fa\u767b\u9646"))),c.a.createElement(g.a,{className:t.paper,elevation:1},c.a.createElement(ue.a,{container:!0,justify:"center"},c.a.createElement(ge.a,{src:this.state.avatar,className:t.bigAvatar})),c.a.createElement(S.a,{variant:"middle",className:t.divider}),c.a.createElement(be.a,{className:t.table},c.a.createElement(Ee.a,null,c.a.createElement(Se.a,null,c.a.createElement(ke.a,{component:"th",scope:"row"},"\u7528\u6237ID"),c.a.createElement(ke.a,{align:"right"},this.state.id)),c.a.createElement(Se.a,null,c.a.createElement(ke.a,{component:"th",scope:"row"},"\u6635\u79f0"),c.a.createElement(ke.a,{align:"right"},this.state.nickname)),c.a.createElement(Se.a,null,c.a.createElement(ke.a,{component:"th",scope:"row"},"\u90ae\u7bb1"),c.a.createElement(ke.a,{align:"right"},this.state.email))))))}}]),t}(c.a.Component),ze=Object(u.withStyles)(function(e){return{root:{flexGrow:1},toolbar:{display:"flex",flexDirection:"row",justifyContent:"space-between",alignItems:"center"},wrapMenu:{display:"flex",flexDirection:"row",justifyContent:"center",alignItems:"center"},menuButton:{marginLeft:-12,marginRight:20},paper:{width:450,margin:"auto",padding:30,marginTop:20,textAlign:"center"},rightIcon:{marginLeft:e.spacing.unit},divider:{margin:"0 0 20px 0"},bigAvatar:{margin:10,width:100,height:100,cursor:"pointer"},fileInput:{display:"none"},b:{display:"flex",width:300}}})(Ie),Ae=a(79),Te=a.n(Ae),Le=function(e){function t(e){var a;return Object(n.a)(this,t),(a=Object(i.a)(this,Object(s.a)(t).call(this,e))).handleInfoClose=function(e,t){if("clickaway"!==t){var n=a.state.info;n.open=!1,a.setState({info:n})}},a.handleKeyDown=function(e){13===e.keyCode&&(a.handleLogin(),e.preventDefault())},a.state={info:{open:!1,variant:"success",message:""},avatar:null,avatarSrc:"http://www.koogo.net:8080/user/statics/avatar/avatar_default.png",email:"",nickname:"",password:"",passwordConfirm:"",nicknameErr:!1,passwordErr:!1,passwordConfirmErr:!1,disabled:!1},a}return Object(o.a)(t,e),Object(r.a)(t,[{key:"componentDidMount",value:function(){var e=this.props.match.params.token.split(".")[1],t=null;try{t=JSON.parse(atob(e))}catch(i){var a=this;return setTimeout(function(){a.props.history.push("/register")},1500),void this.setState({info:{open:!0,variant:"error",message:"\u6ce8\u518c\u9a8c\u8bc1token\u975e\u6cd5!"}})}if(!t.email){var n=this;return setTimeout(function(){n.props.history.push("/register")},1500),void this.setState({info:{open:!0,variant:"error",message:"\u6ce8\u518c\u9a8c\u8bc1token\u975e\u6cd5!"}})}if((new Date).getTime()>1e3*t.exp){var r=this;return setTimeout(function(){r.props.history.push("/register")},1500),void this.setState({info:{open:!0,variant:"error",message:"\u6ce8\u518c\u9a8c\u8bc1token\u5df2\u8fc7\u671f, \u8bf7\u91cd\u65b0\u6ce8\u518c!"}})}console.log(t.email),this.setState({email:t.email,token:this.props.match.params.token})}},{key:"handleAvatar",value:function(e){var t=URL.createObjectURL(e.target.files[0]);this.setState({avatar:e.target.files[0],avatarSrc:t})}},{key:"handleNickname",value:function(e){var t=!1,a=e.target.value;a&&a<4&&(t=!0),this.setState({nickname:a,nicknameErr:t})}},{key:"handlePassword",value:function(e){var t=e.target.value;""!==t?t.length<6?this.setState({password:t,passwordErr:!0}):this.setState({password:t,passwordErr:!1}):this.setState({password:"",passwordErr:!1})}},{key:"handlePasswordConfirm",value:function(e){var t=e.target.value;t===this.state.password?this.setState({passwordConfirm:t,passwordConfirmErr:!1}):this.setState({passwordConfirm:t,passwordConfirmErr:!0})}},{key:"handleComplete",value:function(){var e=this;if(this.state.nickname&&!this.state.nicknameErr)if(this.state.password&&!this.state.passwordErr)if(this.state.passwordConfirmErr)this.setState({info:{open:!0,variant:"error",message:"\u4e24\u6b21\u8f93\u5165\u5bc6\u7801\u4e0d-\u81f4!"}});else{var t=new FormData;t.append("nickname",this.state.nickname),t.append("password",this.state.password),this.state.avatar&&t.append("avatar",this.state.avatar),fetch("http://www.koogo.net:8080/user/auth/complete",{method:"POST",headers:{"x-register-token":this.state.token},body:t}).then(function(e){return 200===e.status?e.json():{code:e.status,message:"internal error"}}).then(function(t){if(200===t.code){var a=e;return setTimeout(function(){a.props.history.push("/profile")},1500),void e.setState({disabled:!1,info:{open:!0,variant:"success",message:"\u6ce8\u518c\u5b8c\u6210, \u8bf7\u767b\u9646!"}})}if(301===t.code){var n=e;return setTimeout(function(){n.props.history.push("/profile")},1500),void e.setState({emailErr:!0,passwordErr:!0,disabled:!1,info:{open:!0,variant:"error",message:"\u8be5\u90ae\u7bb1\u6ce8\u518c\u5df2\u7ecf\u5b8c\u6210, \u8bf7\u767b\u9646!"}})}if(401===t.code){var r=e;return setTimeout(function(){r.props.history.push("/register")},1500),void e.setState({emailErr:!0,disabled:!1,info:{open:!0,variant:"error",message:"\u6ce8\u518c\u9a8c\u8bc1token\u975e\u6cd5\u6216\u5df2\u8fc7\u671f, \u8bf7\u91cd\u65b0\u6ce8\u518c!"}})}e.setState({disabled:!1,info:{open:!0,variant:"error",message:"\u672a\u77e5\u9519\u8bef: "+t.message}})}).catch(function(t){alert(t),e.setState({disabled:!1,info:{open:!0,variant:"error",message:"\u7f51\u7edc\u9519\u8bef: "+t}})}),this.setState({disabled:!0})}else this.setState({info:{open:!0,variant:"error",message:"\u5bc6\u7801\u4e0d\u80fd\u4e3a\u7a7a\u4e14\u957f\u5ea6\u4e0d\u80fd\u5c0f\u4e8e6\u4f4d!"}});else this.setState({info:{open:!0,variant:"error",message:"\u6635\u79f0\u4e0d\u80fd\u4e3a\u7a7a\u4e14\u957f\u5ea6\u4e0d\u80fd\u5c0f\u4e8e4\u4f4d!"}})}},{key:"render",value:function(){var e=this,t=this.props.classes;return c.a.createElement("div",{className:t.container,onSubmit:function(e){alert(e)}},c.a.createElement(b.a,null),c.a.createElement(te.a,{hidden:!this.state.disabled,className:t.progress}),c.a.createElement(Z.a,{anchorOrigin:{vertical:"top",horizontal:"center"},open:this.state.info.open,autoHideDuration:1e3,onClose:this.handleInfoClose},c.a.createElement(Q,{onClose:this.handleInfoClose,variant:this.state.info.variant,message:this.state.info.message})),c.a.createElement(g.a,{className:t.paper},c.a.createElement(k.a,{component:"h1",variant:"h6",color:"primary",className:t.typography},c.a.createElement(N.a,{className:t.cloud}),"\u4e91\u8bb0"),c.a.createElement(S.a,{variant:"middle",className:t.divider}),c.a.createElement("form",{className:t.form,onKeyDown:this.handleKeyDown},c.a.createElement(Te.a,null,c.a.createElement("input",{accept:"image/*",className:t.fileInput,id:"avatar-file",multiple:!0,type:"file",onChange:function(t){return e.handleAvatar(t)}}),c.a.createElement("label",{htmlFor:"avatar-file",className:t.avatarLabel},c.a.createElement(ge.a,{alt:"",src:this.state.avatarSrc,className:t.bigAvatar}),c.a.createElement(ne.a,{variant:"body2"},"\u70b9\u51fb\u4e0a\u4f20"))),c.a.createElement(O.a,{className:t.textField,id:"complete-nickname",variant:"filled",label:"\u6635\u79f0",autoComplete:"off",defaultValue:this.state.nickname,error:this.state.nicknameErr,disabled:this.state.disabled,onChange:function(t){return e.handleNickname(t)},onBlur:function(t){return e.handleNickname(t)},required:!0}),c.a.createElement(O.a,{className:t.textField,variant:"filled",label:"\u90ae\u7bb1",autoComplete:"off",value:this.state.email,disabled:!0}),c.a.createElement(O.a,{className:t.textField,id:"complete-password",variant:"filled",label:"\u5bc6\u7801",autoComplete:"new-password",type:"password",value:this.state.password,onChange:function(t){return e.handlePassword(t)},onBlur:function(t){return e.handlePassword(t)},error:this.state.passwordErr,disabled:this.state.disabled,required:!0}),c.a.createElement(O.a,{className:t.textField,id:"complete-password-confirm",variant:"filled",label:"\u5bc6\u7801\u786e\u8ba4",autoComplete:"new-password",type:"password",value:this.state.passwordConfirm,onChange:function(t){return e.handlePasswordConfirm(t)},onBlur:function(t){return e.handlePasswordConfirm(t)},error:this.state.passwordConfirmErr,disabled:this.state.disabled,required:!0}),c.a.createElement(E.a,{variant:"outlined",color:"primary",className:t.button,onClick:function(){return e.handleComplete()},disabled:this.state.disabled},"\u5b8c\u6210\u6ce8\u518c"))))}}]),t}(c.a.Component),Be=Object(u.withStyles)(function(e){return{container:{height:"100vh",background:"linear-gradient(to right, #348AC7, #7474BF);",overflow:"hidden"},progress:{position:"fixed",width:"100%",background:"transparent"},paper:{maxWidth:400,margin:"50px auto",padding:"20px 0 40px 0"},typography:{textAlign:"center",verticalAlign:"middle",fontSize:"18px",lineHeight:"40px"},cloud:{verticalAlign:"top",fontSize:"35px",marginRight:10},divider:{margin:"0 0 20px 0"},form:{display:"flex",flexDirection:"column",justifyContent:"center",alignItems:"center"},textField:{width:"70%",background:"transparent"},avatarLabel:{textAlign:"center",cursor:"pointer",marginBottom:10},button:{fontSize:"18px",height:45,width:"40%",marginTop:30,marginBottom:5},bigAvatar:{width:100,height:100,cursor:"pointer"},fileInput:{display:"none"}}})(Le),Fe=function(e){function t(){return Object(n.a)(this,t),Object(i.a)(this,Object(s.a)(t).apply(this,arguments))}return Object(o.a)(t,e),Object(r.a)(t,[{key:"render",value:function(){return c.a.createElement(h.a,null,c.a.createElement(p.a,{path:"/",exact:!0,component:ze}),c.a.createElement(p.a,{path:"/Register",component:ie}),c.a.createElement(p.a,{path:"/Complete/:token",component:Be}),c.a.createElement(p.a,{path:"/login",component:le}),c.a.createElement(p.a,{path:"/profile",component:ze}))}}]),t}(c.a.Component);m.a.render(c.a.createElement(Fe,null),document.querySelector("#root"))}},[[131,1,2]]]);
//# sourceMappingURL=main.a7bd3a66.chunk.js.map