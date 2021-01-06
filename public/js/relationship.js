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
  
  async function getAccessToken(jwt){
  await superagent
    .get('/api/v1/auth/ics/refresh_token')
    .set('accept', 'json')
    .then(function (res) {
        console.log(res.body);
        jwt.token = res.body.token;
        jwt.token_expiry = res.body.token_expiry;
        jwt.token_type = res.body.type;
        setSession(jwt);
    })
    .catch(function (err) {
      //logout if err occur
      if (err.response) {
        console.log(err);
      }
      let url = window.location.href.split('/');
      if(url.length > 4){
        logout();
      }
    });
  }
  
  function logout(){
    jwt.token_expiry = new Date('-1');
    jwt.token = '';
    setSession(jwt);
  
    //ask server to set refresh token invalid
    var request = new XMLHttpRequest(); 
    request.open('GET', '/api/v1/auth/ics/logout', true);
    request.send();
    request.onerror = showMessage();
    window.location.replace("/");
  }
  
  function setJWTFromSession(jwt){
    jwt.token = Session.get('token');
    jwt.token_expiry = new Date(Session.get('token_expiry'));
    jwt.token_type = Session.get('token_type');
  }
  
  function setSession(jwt){
    Session.set('token', jwt.token);
    Session.set('token_expiry', jwt.token_expiry);
    Session.set('token_type', jwt.type);
  }
  
  var jwt = {
    token: Session.get('token'),
    token_expiry: new Date(Session.get('token_expiry')),
    token_type: Session.get('token_type'),
  };
  
  async function checkToken(jwt){
    await getAccessToken(jwt);
  }
  
  //check if user login
  function isLogin(){
    //1. if token is expire, try to get new access token
    if(jwt.token_expiry < new Date() && jwt.token != '' && typeof(jwt.token) != 'undefined'){
      checkToken();
    }
    //2. check if token is invalid
    if(jwt.token_expiry < new Date() || jwt.token == '' || typeof(jwt.token) == 'undefined'){
      return false
    }
    //3. token is valid
    document.querySelector('#login').style.display = 'none';
    document.querySelector('#signup').style.display = 'none';
    document.querySelector('#logout').style.display = 'block';

    let tokenString = Session.get('token_type') + ' ' + Session.get('token');
    let url = window.location.href.split('/');
    let id = url[url.length-1];

    document.querySelector('#profile').href = '/profile/' + id;
    document.querySelector('#history').href = '/history/' + id;
    document.querySelector('#myProfile').href = '/profile/' + id;
    document.querySelector('#myHistory').href = '/history/' + id;
    document.querySelector('#myFriends').href = '/friends/' + id;
    document.querySelector('#myFollowers').href = '/followers/' + id;


    return true;
  }
  isLogin();