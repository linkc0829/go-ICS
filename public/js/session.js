/**
 * Implements cookie-less JavaScript session variables
 * v1.0
 *
 * By Craig Buckler, Optimalworks.net
 * http://dreamerslab.com/blog/tw/javascript-session/
 *
 */
  
 if (JSON && JSON.stringify && JSON.parse) var Session = Session || (function() {
  
  // cache window 物件
  var win = window.top || window;
   
  // 將資料都存入 window.name 這個 property
  var store = (win.name ? JSON.parse(win.name) : {});
   
  // 將要存入的資料轉成 json 格式
  function Save() {
    win.name = JSON.stringify(store);
  };
   
  // 在頁面 unload 的時候將資料存入 window.name
  if (window.addEventListener) window.addEventListener("unload", Save, false);
  else if (window.attachEvent) window.attachEvent("onunload", Save);
  else window.onunload = Save;
 
  // public methods
  return {
   
    // 設定一個 session 變數
    set: function(name, value) {
      store[name] = value;
    },
     
    // 列出指定的 session 資料
    get: function(name) {
      return (store[name] ? store[name] : undefined);
    },
     
    // 清除資料 ( session )
    clear: function() { store = {}; },
     
    // 列出所有存入的資料
    dump: function() { return JSON.stringify(store); }
  
  };
  
 })();