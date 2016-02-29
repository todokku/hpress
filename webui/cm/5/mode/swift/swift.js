!function(e){"object"==typeof exports&&"object"==typeof module?e(require("../../lib/codemirror")):"function"==typeof define&&define.amd?define(["../../lib/codemirror"],e):e(CodeMirror)}(function(e){"use strict";function t(e){for(var t={},n=0;n<e.length;n++)t[e[n]]=!0;return t}function n(e,t,n){if(e.eatSpace())return null;var r=e.peek();if("/"==r){if(e.match("//"))return e.skipToEnd(),"comment";if(e.match("/*"))return t.tokenize.push(o),o(e,t);if(e.match(h))return"string-2"}if(p.indexOf(r)>-1)return e.next(),"operator";if(l.indexOf(r)>-1)return e.match(s),"punctuation";if('"'==r||"'"==r){e.next();var v=i(r);return t.tokenize.push(v),v(e,t)}if(e.match(d))return"number";if(e.match(_))return"property";if(e.match(m)){var k=e.current();return u.hasOwnProperty(k)?(a.hasOwnProperty(k)&&(t.prev="define"),"keyword"):c.hasOwnProperty(k)?"variable-2":f.hasOwnProperty(k)?"atom":"define"==n?"def":"variable"}return e.next(),null}function r(){var e=0;return function(t,r,i){var o=n(t,r,i);if("punctuation"==o)if("("==t.current())++e;else if(")"==t.current()){if(0==e)return t.backUp(1),r.tokenize.pop(),r.tokenize[r.tokenize.length-1](t,r);--e}return o}}function i(e){return function(t,n){for(var i,o=!1;i=t.next();)if(o){if("("==i)return n.tokenize.push(r()),"string";o=!1}else{if(i==e)break;o="\\"==i}return n.tokenize.pop(),"string"}}function o(e,t){return e.match(/^(?:[^*]|\*(?!\/))*/),e.match("*/")&&t.tokenize.pop(),"comment"}var u=t(["var","let","class","deinit","enum","extension","func","import","init","protocol","static","struct","subscript","typealias","as","dynamicType","is","new","super","self","Self","Type","__COLUMN__","__FILE__","__FUNCTION__","__LINE__","break","case","continue","default","do","else","fallthrough","if","in","for","return","switch","where","while","associativity","didSet","get","infix","inout","left","mutating","none","nonmutating","operator","override","postfix","precedence","prefix","right","set","unowned","weak","willSet"]),a=t(["var","let","class","enum","extension","func","import","protocol","struct","typealias","dynamicType","for"]),f=t(["Infinity","NaN","undefined","null","true","false","on","off","yes","no","nil","null","this","super"]),c=t(["String","bool","int","string","double","Double","Int","Float","float","public","private","extension"]),p="+-/*%=|&<>#",l=";,.(){}[]",s=/^(?:[()\[\]{},:`=;]|\.\.?\.?)/,d=/^-?(?:(?:[\d_]+\.[_\d]*|\.[_\d]+|0o[0-7_\.]+|0b[01_\.]+)(?:e-?[\d_]+)?|0x[\d_a-f\.]+(?:p-?[\d_]+)?)/i,m=/^[_A-Za-z$][_A-Za-z$0-9]*/,_=/^[@\.][_A-Za-z$][_A-Za-z$0-9]*/,h=/^\/(?!\s)(?:\/\/)?(?:\\.|[^\/])+\//;e.defineMode("swift",function(){return{startState:function(){return{prev:null,tokenize:[]}},token:function(e,t){var r=t.prev;t.prev=null;var i=t.tokenize[t.tokenize.length-1]||n,o=i(e,t,r);return o&&"comment"!=o?t.prev||(t.prev=o):t.prev=r,o},lineComment:"//",blockCommentStart:"/*",blockCommentEnd:"*/"}}),e.defineMIME("text/x-swift","swift")});