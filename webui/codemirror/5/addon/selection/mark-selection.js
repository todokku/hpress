// CodeMirror, copyright (c) by Marijn Haverbeke and others
// Distributed under an MIT license: http://codemirror.net/LICENSE
// Because sometimes you need to mark the selected *text*.
//
// Adds an option 'styleSelectedText' which, when enabled, gives
// selected text the CSS class given as option value, or
// "CodeMirror-selectedtext" when the value is not a string.
(function(e){typeof exports=="object"&&typeof module=="object"?e(require("../../lib/codemirror")):typeof define=="function"&&define.amd?define(["../../lib/codemirror"],e):e(CodeMirror)})(function(e){"use strict";function t(e){e.operation(function(){f(e)})}function n(e){e.state.markedSelection.length&&e.operation(function(){u(e)})}function o(e,t,n,o){if(s(t,n)==0)return;var u=e.state.markedSelection,a=e.state.markedSelectionStyle;for(var f=t.line;;){var l=f==t.line?t:i(f,0),c=f+r,h=c>=n.line,p=h?n:i(c,0),d=e.markText(l,p,{className:a});o==null?u.push(d):u.splice(o++,0,d);if(h)break;f=c}}function u(e){var t=e.state.markedSelection;for(var n=0;n<t.length;++n)t[n].clear();t.length=0}function a(e){u(e);var t=e.listSelections();for(var n=0;n<t.length;n++)o(e,t[n].from(),t[n].to())}function f(e){if(!e.somethingSelected())return u(e);if(e.listSelections().length>1)return a(e);var t=e.getCursor("start"),n=e.getCursor("end"),i=e.state.markedSelection;if(!i.length)return o(e,t,n);var f=i[0].find(),l=i[i.length-1].find();if(!f||!l||n.line-t.line<r||s(t,l.to)>=0||s(n,f.from)<=0)return a(e);while(s(t,f.from)>0)i.shift().clear(),f=i[0].find();s(t,f.from)<0&&(f.to.line-t.line<r?(i.shift().clear(),o(e,t,f.to,0)):o(e,t,f.from,0));while(s(n,l.to)<0)i.pop().clear(),l=i[i.length-1].find();s(n,l.to)>0&&(n.line-l.from.line<r?(i.pop().clear(),o(e,l.from,n)):o(e,l.to,n))}e.defineOption("styleSelectedText",!1,function(r,i,s){var o=s&&s!=e.Init;i&&!o?(r.state.markedSelection=[],r.state.markedSelectionStyle=typeof i=="string"?i:"CodeMirror-selectedtext",a(r),r.on("cursorActivity",t),r.on("change",n)):!i&&o&&(r.off("cursorActivity",t),r.off("change",n),u(r),r.state.markedSelection=r.state.markedSelectionStyle=null)});var r=8,i=e.Pos,s=e.cmpPos});