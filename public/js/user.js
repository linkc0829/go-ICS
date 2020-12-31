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

function getAccessToken(){
  const res = superagent
    .get('/api/v1/auth/ics/refresh_token')
    .set('accept', 'json')
    .then(function (res) {
        clearMessage();
        token_expiry = Date.parse(res.body.token_expiry);
        Session.set('token', res.body.token);
        Session.set('token_expiry', token_expiry);
        token = Session.get('token')
        token_expiry = Session.get('token_expiry')
        input.value = '';
    })
    .catch(function (err) {
        if (err.response) {
            showMessage(err.response.message);
        }
    });
}

document.addEventListener('DOMContentLoaded', function () {
    /* Event listeners related to TODO creation. */
    (function () {
        /* The input entry listens to ENTER press event. */
        let form = document.querySelector('#add');

        let id = document.querySelector('#id');
        let userId = document.querySelector('#userId');
        let email = document.querySelector('#email');
        let nickName = document.querySelector('#nickName');
        let cost = document.querySelector('#cost');
        let income = document.querySelector('#income');


        id.addEventListener('keydown', function (ev) {
            if (ev.key === 13) {
                ev.preventDefault();
                if(cost.checked){
                    createCost();
                }
                else if(income){
                    createIncome();
                }
            }
        }, false);

        /* The `Add` button listens to click event. */
        let btn = form.querySelector('button');

        btn.addEventListener('click', function (ev) {
            ev.preventDefault();
            if(cost.checked){
                createCost();
            }
            else if(income){
                createIncome();
            }
        }, false);

        function createIncome() {
            let target = `/api/v1/user` + id.value;
            let data = {}
            if(userId.value) {
                data.userId = userId.value
            }
            if(email.value) {
                data.email = email.value
            }
            if(nickName.value) {
                data.nickName = nickName.value
            }
            console.log(data)
                
            superagent
                .post(target)
                .send(data)
                .set('accept', 'json')
                .set('Authorization', token)
                .then(function (res) {
                    console.log('finished');
                    clearMessage();
                    let result = JSON.stringify(res.body);       

                    let resp = document.getElementById('response');
                    resp.innerHTML = result;

                    id.value = res.body.UpdateUser.Id;
                    userId.value = res.body.UpdateUser.UserId;
                    email.value = res.body.UpdateUser.Email;
                    nickName.value = res.body.UpdateUser.NickName;

                })
                .catch(function (err) {
                    if (err.response) {
                        showMessage(err.response.message);
                    }
                });
        }
        
    })();
})

function showMessage(msg) {
    let div = document.createElement('div');

    div.classList.add('alert');
    div.classList.add('alert-warning');
    div.setAttribute('role', 'alert');

    div.innerText = msg;

    let message = document.getElementById('message');

    message.innerHTML = '';
    message.appendChild(div);
}

function clearMessage() {
    let message = document.getElementById('message');

    message.innerHTML = '';
}

var token = Session.get('token');
var token_expiry = Session.get('token_expiry');
var now = new Date()

if (typeof(token_expiry) == 'undefined' || token_expiry < now){
  getAccessToken();
}