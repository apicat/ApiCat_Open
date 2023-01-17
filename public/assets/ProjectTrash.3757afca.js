import{E as D,S as F,a as x,b as T,e as B,c as L,K as V,L as P,H as I,v as N}from"./element-plus.aa12b6bd.js";import{ar as p,o as f,Y as C,Z as r,$ as a,aA as $,ak as g,K as U,a$ as A,r as h,x as K,h as O,T as R,j as b,W as M,bn as w}from"./vendor.b7ec7f69.js";import{c as q,N as k,O as E,Q as H,B as W,U as y,V as Q,W as Y}from"./index.6e925652.js";const Z={emits:["on-ok"],data(){return{project_id:this.$route.params.project_id||"",isShow:!1,isLoading:!1,dir:[],form:{dir_id:[]},rules:{dir_id:{required:!0,min_len:1,message:"\u8BF7\u9009\u62E9\u5206\u7C7B",trigger:"change",type:"array"}}}},watch:{isShow:function(){!this.isShow&&this.reset()},"form.dir_id":function(){let t=this.form.dir_id.slice(-1)[0];t!==void 0&&(this.document.node_id=t)}},methods:{show(t){if(!t)throw new Error("\u6587\u6863\u4FE1\u606F\u4E0D\u80FD\u4E3A\u7A7A\uFF01");this.document=t,this.isShow=!0},hide(){this.isShow=!1},onCloseBtnClick(){this.isShow=!1,this.reset()},reset(){this.$refs.teamForm.resetFields()},handleSubmit(t){this.$refs[t].validate(e=>e&&this.submit())},submit(){this.isLoading=!0,k(this.document).then(t=>{this.onCloseBtnClick(),D({type:"success",closable:!0,message:p("span",null,["\u6587\u6863\u6062\u590D\u6210\u529F\uFF0C",p("a",{class:"text-blue-600",href:E(this.project_id,this.document.doc_id,!1)},"\u67E5\u770B\u8BE6\u60C5")])}),this.$emit("on-ok")}).finally(()=>{this.isLoading=!1})},transferDir(t){let e=[],d=(c,o)=>{(c||[]).forEach(s=>{let _={value:s.id,label:s.title,children:[]};o.push(_),s.sub_nodes&&s.sub_nodes.length&&d(s.sub_nodes,_.children)})};return d(t,e),[{value:0,label:"\u6839\u76EE\u5F55"}].concat(e)},async getDocumentDirList(t){H(t).then(({data:e})=>{this.dir=this.transferDir(e)})}},mounted(){this.getDocumentDirList(this.project_id)}},z=g("\u6682\u65E0\u6570\u636E"),G=g(" \u53D6\u6D88 "),J=g(" \u786E\u5B9A ");function X(t,e,d,c,o,s){const _=F,m=x,n=T,l=B,i=L;return f(),C(i,{modelValue:o.isShow,"onUpdate:modelValue":e[4]||(e[4]=u=>o.isShow=u),width:400,class:"show-footer-line vertical-center-modal","close-on-click-modal":!1,"append-to-body":"",title:"\u539F\u6587\u6863\u6240\u5728\u5206\u7C7B\u5DF2\u88AB\u5220\u9664\uFF0C\u8BF7\u9009\u62E9\u5176\u4ED6\u5206\u7C7B"},{footer:r(()=>[a(l,{onClick:e[2]||(e[2]=u=>s.onCloseBtnClick())},{default:r(()=>[G]),_:1}),a(l,{loading:o.isLoading,type:"primary",onClick:e[3]||(e[3]=u=>s.handleSubmit("teamForm"))},{default:r(()=>[J]),_:1},8,["loading"])]),default:r(()=>[a(n,{ref:"teamForm",model:o.form,rules:o.rules,"label-position":"top",style:{"margin-bottom":"-19px"},onKeyup:e[1]||(e[1]=$(u=>s.handleSubmit("teamForm"),["enter"]))},{default:r(()=>[a(m,{label:"",prop:"dir_id",class:"hide_required"},{default:r(()=>[a(_,{modelValue:o.form.dir_id,"onUpdate:modelValue":e[0]||(e[0]=u=>o.form.dir_id=u),class:"w-full",options:o.dir,props:{checkStrictly:!0},placeholder:"\u8BF7\u9009\u62E9\u5206\u7C7B"},{empty:r(()=>[z]),_:1},8,["modelValue","options"])]),_:1})]),_:1},8,["model","rules"])]),_:1},8,["modelValue"])}var ee=q(Z,[["render",X]]);const te=b("span",null,"\u56DE\u6536\u7AD9",-1),oe=["href"],se=["onClick"],ne=U({setup(t){const e=W(),{projectInfo:d}=A(e),c=h(),o=h([]),s=h(!1),_=n=>{w.start();const l={project_id:e.projectInfo.id,doc_id:n.id};k(l).then(({status:i})=>{if(i===y.NO_PARENT_DIR){c.value&&c.value.show(l);return}i===y.OK&&(D({type:"success",showClose:!0,message:()=>p("span",null,["\u6587\u6863\u6062\u590D\u6210\u529F\uFF0C",p("a",{class:"text-blue-600",href:E(e.projectInfo.id,n.id,!1)},"\u67E5\u770B\u8BE6\u60C5")])}),m())}).catch(i=>{}).finally(()=>{w.done()})},m=async()=>{s.value=!0;try{const{data:n}=await Q(e.projectInfo.id);o.value=(n||[]).map(l=>(l.previewUrl=Y({project_id:e.projectInfo.id,doc_id:l.id}),l))}catch{}finally{s.value=!1}};return K(()=>d.value,async()=>{d.value&&d.value.id&&await m()},{immediate:!0}),(n,l)=>{const i=V,u=P,j=I,S=N;return f(),O(M,null,[R((f(),C(j,{shadow:"never","body-style":{padding:0}},{header:r(()=>[te]),default:r(()=>[a(u,{data:o.value,"empty-text":"\u6682\u65E0\u6570\u636E"},{default:r(()=>[a(i,{prop:"title",label:"\u6587\u6863\u540D\u79F0","show-overflow-tooltip":""}),a(i,{prop:"deleted_at",label:"\u5220\u9664\u65F6\u95F4"}),a(i,{prop:"remaining",label:"\u5269\u4F59"}),a(i,{label:"\u64CD\u4F5C"},{default:r(({row:v})=>[b("a",{class:"mr-3 text-blue-600 cursor-pointer",target:"_blank",href:v.previewUrl},"\u9884\u89C8",8,oe),b("span",{class:"mr-3 text-blue-600 cursor-pointer",href:"javascript:void(0)",onClick:re=>_(v)},"\u6062\u590D",8,se)]),_:1})]),_:1},8,["data"])]),_:1})),[[S,s.value]]),a(ee,{ref_key:"restoreDocumentModal",ref:c,onOnOk:m},null,512)],64)}}});export{ne as default};
