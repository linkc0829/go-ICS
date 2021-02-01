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

async function initProfile(jwt){

  await silentRefresh(jwt);

  let tokenString = Session.get('token_type') + ' ' + Session.get('token');
  let url = window.location.href.split('/');
  let id = url[url.length-1];
  let query = '{\
      me{\
        id\
        friends{\
          id\
        }\
      }\
    }';

  await superagent
  .post('/api/v1/graph')
  .set('accept', 'json')
  .set('Authorization', tokenString)
  .send({'Query': query,})
  .then(function (res) {
      let out = res.body.data.me.friends;
      //init addFriend button
      if(res.body.data.me.id != id){
        document.querySelector('#addFriend').style.display = 'block';
        document.querySelector('#unfriend').style.display = 'none';
      }
      //check if being friend
      for(let i = 0; i < out.length; i++){
        if(out[i].id == id){
          document.querySelector('#addFriend').style.display = 'none';
          document.querySelector('#unfriend').style.display = 'block';
        }
      }
      //watching my own profile, hide friend sector
      if(res.body.data.me.id == id || id == ''){
        document.querySelector('#addFriend').style.display = 'none';
        document.querySelector('#unfriend').style.display = 'none';
      }
      //watching others profile, hide add sector, remove eventlistners, update&delete button
      if(res.body.data.me.id != id){
        document.querySelector('#add').style.display = 'none';
      }
      //add my account links
      let myID = res.body.data.me.id;
      Session.set('myID', myID);
      document.querySelector('#myProfile').href = '/profile/' + myID;
      document.querySelector('#myHistory').href = '/history/' + myID;
      document.querySelector('#myFriends').href = '/friends/' + myID;
      document.querySelector('#myFollowers').href = '/followers/' + myID;

      if(url[3] == 'profile' || url[3] == 'history'){
        getUserProfile(currentUser);
      }

      let date = new Date();
      let date_month = date.getMonth()+1;
      let occurDate = date.getFullYear() + '-' + date_month + '-' + date.getDate();

      if(url[3] != 'history'){
        document.querySelector('#occurDate').min = occurDate;
      }

      document.querySelector('#login').style.display = 'none';
      document.querySelector('#signup').style.display = 'none';
      document.querySelector('#logout').style.display = 'block';

  })
  .catch(function (err) {
      alert(err);
  });
}

//get UserId through gql API
async function getUserProfile(currentUser){
  let tokenString = Session.get('token_type') + ' ' + Session.get('token');
  let url = window.location.href.split('/');
  let id = url[url.length-1];
  
  let query = '{\
      getUser(id: "' + id + '"){\
        id\
        userId\
        email\
        nickName\
        friends{\
          id\
        }\
        followers{\
          id\
        }\
        role\
      }\
    }';

  await superagent
  .post('/api/v1/graph')
  .set('accept', 'json')
  .set('Authorization', tokenString)
  .send({'Query': query,})
  .then(function (res) {
      
      let out = res.body.data.getUser;
      currentUser.id = out.id;
      currentUser.userId = out.userId;
      currentUser.nickName = out.nickName;
      currentUser.email = out.email;
      currentUser.role = out.role;
      
      currentUser.friends = (out.friends == null ? []: [...out.friends]);
      currentUser.followers = (out.followers == null ? []: [...out.followers]);

      let resp = document.getElementById('response');
      resp.innerHTML = JSON.stringify(res.body);

      document.querySelector('#profile').href = '/profile/' + out.id;
      document.querySelector('#history').href = '/history/' + out.id;

      if(url[3] == 'history'){
        initHistory(currentUser, INCOME, 30);
        initHistory(currentUser, COST, 30);
      }
      else{
        initPortfolio(currentUser, INCOME);
        initPortfolio(currentUser, COST);
      }
  })
  .catch(function (err) {
    alert(err);
  });
}

async function silentRefresh(jwt){
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
    alert(err + " get access token failed, now logout");
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
  Session.set('token_type', jwt.token_type);
}

function disableUpdate(type){
  window.loadPortfolioItem = function(){};
  let parent = (type == INCOME? document.querySelector('#income-list'): document.querySelector('#cost-list'));
  let children = parent.children;
  for(let i = 0; i < children.length; i++){
    children[i].querySelector('#delete').remove();
    children[i].querySelector('#update').remove();
    children[i].querySelector('#vote').remove();
  }
}

const COST = 'COST';
const INCOME = 'INCOME';

var jwt = {
  token: Session.get('token'),
  token_expiry: new Date(Session.get('token_expiry')),
  token_type: Session.get('token_type'),
};

async function checkToken(jwt){
  await silentRefresh(jwt);
}

//check if user login
function isLogin(){
  console.log("isLogin")
  //1. if token is expire, try to get new access token
  if(jwt.token_expiry < new Date() && jwt.token != '' && typeof(jwt.token) != 'undefined'){
    checkToken();
  }
  //2. if token is invalid
  if(jwt.token_expiry < new Date() || jwt.token == '' || typeof(jwt.token) == 'undefined'){
    document.querySelector('#login').style.display = 'block';
    document.querySelector('#signup').style.display = 'block';
    document.querySelector('#history').style.display = 'none';
    document.querySelector('#logout').style.display = 'none';
    return false
  }
  //3. token is valid
  document.querySelector('#login').style.display = 'none';
  document.querySelector('#signup').style.display = 'none';
  document.querySelector('#logout').style.display = 'block';

  //set myaccount links
  let myID = Session.get('myID');
  document.querySelector('#myProfile').href = '/profile/' + myID;
  document.querySelector('#myHistory').href = '/history/' + myID;
  document.querySelector('#myFriends').href = '/friends/' + myID;
  document.querySelector('#myFollowers').href = '/followers/' + myID;
  
  return true;
}

var currentUser = {};
var url = window.location.href.split('/');
var last = url[url.length-1];
//when open new window, try to get new access token
if(typeof(jwt.token) == 'undefined' && typeof(jwt.token_type) == 'undefined'){
  initProfile(jwt);
}
else {
  isLogin();
  if(last.length == 24){
    getUserProfile(currentUser);
  }
}

if(last.length != 24){
  document.querySelector('#profile').style.display = 'none';
  document.querySelector('#history').style.display = 'none';
  document.querySelector('#addFriend').style.display = 'none';
}
if(url[3] == 'history'){
  disableUpdate(INCOME);
  disableUpdate(COST);
}

//load user portfolio, for init or reload
async function initHistory(user, type, range){
  let target = '/api/v1/user/' + user.id + (type==INCOME? '/income': '/cost') + '/history?range=' + range;
  let portfolioList = (type==INCOME? document.querySelector('#income-list'):document.querySelector('#cost-list'));
  portfolioList.innerHTML = '';
  await superagent
      .get(target)
      .set('accept', 'json')
      .set('Authorization', jwt.token)
      .then(function (res) {
          console.log(res);
          let portfolio;
          if(type == INCOME){
            portfolio = res.body.GetUserIncomeHistory;
          }
          else{
            portfolio = res.body.GetUserCostHistory;
          }
          for (let i = 0; i < portfolio.length; i++) {
            addPortfolio(portfolio[i], type);
        }
      })
      .catch(function (err) {
        alert(err);
      });
}

function reloadHistory(range){
  isLogin();
  initHistory(currentUser, INCOME, range);
  initHistory(currentUser, COST, range);
  disableUpdate();
}

function casePortfolioType(type, upper){
  if(type == COST){
    return upper? 'Cost':'cost';
  }
  else{
    return upper? 'Income':'income';
  }
}

function createPortfolioRecord(type){

  let ret = "<div class='form-row'>\
  <div class='col-md-6 border'>\
  <label for='description' id='description_label' style='width:100%;'>" + casePortfolioType(type, true) + " Description</label>\
  <input type='text' class='form-control' name='description' placeholder='" + casePortfolioType(type, true) + " Description' id='description_input' style='display: none;'>\
  </div>\
  <div class='col-md-2 border'>\
  <label for='amount' id='amount_label' style='width:100%;'>Amount: $888</label>\
  <input type='text' class='form-control' name='amount' placeholder='888' id='amount_input' style='display: none;'>\
  </div>\
  <div class='col-md-2 border'>\
  <label for='privacy' id='privacy_label' style='width:100%;'>FRIEND</label>\
  <select name='privacy' id='privacy_input' class='custom-select' style='display: none; font-size: 0.9rem; margin-top: 5px;'>\
    <option value='PRIVATE' selected>PRIVATE</option>\
    <option value='FRIEND'>FRIEND</option>\
    <option value='PUBLIC'>PUBLIC</option>\
  </select>\
  </div>\
  <div class='col-md-2'>\
  <button class='btn btn-secondary' id='update'>Update</button>\
  </div>\
  </div>\
  <div class='form-row'>\
  <div class='col-md-3 border'>"

  if(type == INCOME){
    ret += "<label for='category' id='category_label' style='width:100%;'>INVESTMENT</label>\
    <select name='category' id='category_input' class='custom-select' style='display: none;'>\
      <option selected>Choose a category</option>\
      <option value='INVESTMENT'>INVESTMENT</option>\
      <option value='SALARY'>SALARY</option>\
      <option value='PARTTIME'>PART TIME</option>\
      <option value='OTHERS'>OTHERS</option>\
    </select>"
  }
  else{
    ret += "<label for='category' id='category_label' style='width:100%;'>LEARNING</label>\
    <select name='category' id='category_input' style='display: none;' class='custom-select'>\
      <option selected>Choose a category</option>\
      <option value='INVESTMENT'>INVESTMENT</option>\
      <option value='DAILY'>DAILY</option>\
      <option value='LEARNING'>LEARNING</option>\
      <option value='CHARITY'>CHARITY</option>\
      <option value='OTHERS'>OTHERS</option>\
    </select>"
  }
  ret +="</div>\
  <div class='col-md-3 border'>\
  <label for='occurDate' id='occurDate_label' style='width:100%;'>@ 2021-01-01</label>\
  <input type='date' class='form-control' name='occurDate' placeholder='occurDate' min='' value='2021-01-01' style='display: none;' id='occurDate_input'>\
  </div>\
  <div class='col-md-2 border'>\
  <label id='vote_label'>Vote: 555</label>\
  </div>\
  <div class='col-md-2'>\
  <button class='btn btn-primary' id='vote'>Vote</button>\
  </div>\
  <div class='col-md-2'>\
  <button class='btn btn-danger' id='delete'>Delete</button>\
  </div>\
  </div>"
  return ret;
}

//load user portfolio, for init or reload
async function initPortfolio(user, type){
  let target = '/api/v1/user/' + user.id + (type==INCOME? '/income': '/cost') ;
  let portfolioList = (type==INCOME? document.querySelector('#income-list'):document.querySelector('#cost-list'));
  portfolioList.innerHTML = '';
  await superagent
      .get(target)
      .set('accept', 'json')
      .set('Authorization', jwt.token)
      .then(function (res) {
        let portfolio;
        if(type == INCOME){
          portfolio = res.body.GetUserIncome;
        }
        else{
          portfolio = res.body.GetUserCost;
        }
        for (let i = 0; i < portfolio.length; i++) {
          addPortfolio(portfolio[i], type);
        }
      })
      .catch(function (err) {
          alert(err)
      });
}

/* Add a Portfolio item. */
function addPortfolio(res, type) {
  if(type != COST && type != INCOME) return;
  let id = res.Id;
  let description = res.Description;
  let amount = res.Amount;
  let category = res.Category;
  let date = new Date(res.OccurDate);
  let date_month = date.getMonth()+1;
  let occurDate = date.getFullYear() + '-' + date_month + '-' + date.getDate();
  let vote = (res.Vote? res.Vote.length:0);
  let privacy = res.Privacy;

  let form = document.createElement('form');
  form.setAttribute('id', id);
  form.innerHTML = createPortfolioRecord(type)
  let portfolioList = document.querySelector((type==INCOME? '#income-list': '#cost-list'));
  portfolioList.appendChild(form);
 
  let description_label = form.querySelector('#description_label');
  description_label.innerText = description;
  description_label.htmlFor = 'description_' + id;
  description_label.addEventListener('click', ()=>{ loadPortfolioItem(id, 'description', type)})
  let description_input = form.querySelector('#description_input');
  description_input.name = 'description_' + id;
  description_input.value = description;

  let amount_label = form.querySelector('#amount_label');
  amount_label.innerText = amount;
  amount_label.htmlFor = 'amount_' + id;
  amount_label.addEventListener('click', ()=>{ loadPortfolioItem(id, 'amount', type)})
  let amount_input = form.querySelector('#amount_input');
  amount_input.name = 'amount_' + id;
  amount_input.value = amount;

  let privacy_label = form.querySelector('#privacy_label');
  privacy_label.innerText = privacy;
  privacy_label.htmlFor = 'privacy_' + id;
  privacy_label.addEventListener('click', ()=>{ loadPortfolioItem(id, 'privacy', type)})
  let privacy_input = form.querySelector('#privacy_input');
  privacy_input.name = 'privacy_' + id;
  privacy_input.value = privacy;

  let category_label = form.querySelector('#category_label');
  category_label.innerText = category;
  category_label.htmlFor = 'category_' + id;
  category_label.addEventListener('click', ()=>{ loadPortfolioItem(id, 'category', type)})
  let category_input = form.querySelector('#category_input');
  category_input.name = 'category_' + id;
  category_input.value = category;

  let occurDate_label = form.querySelector('#occurDate_label');
  occurDate_label.innerText = occurDate;
  occurDate_label.htmlFor = 'occueDate_' + id;
  occurDate_label.addEventListener('click', ()=>{ loadPortfolioItem(id, 'occurDate', type)})
  let occurDate_input = form.querySelector('#occurDate_input');
  occurDate_input.name = 'occueDate_' + id;
  occurDate_input.value = occurDate;

  let vote_label = form.querySelector('#vote_label');
  vote_label.innerText = vote;
  let voteBtn = form.querySelector('#vote');
  voteBtn.addEventListener('click', (e)=>{
    isLogin();
    e.preventDefault();
    let target = '/api/v1/' + casePortfolioType(type, false) + '/vote/' + id;
    superagent
        .put(target)
        .set('accept', 'json')
        .set('Authorization', jwt.token)
        .then(function (res) {
          if(res.errors){
            alert(res.errors);
            initPortfolio(currentUser, type);
          }
          vote_label.innerText = (type==INCOME? res.body.VoteIncome:res.body.VoteCost);
        })
        .catch(function (err) {
          alert(err);
          initPortfolio(currentUser, type);
        });
  })

  let deleteBtn = form.querySelector('#delete');
  deleteBtn.addEventListener('click', (e)=>{
    isLogin();
    e.preventDefault();
    form.parentNode.removeChild(form);
    let target = '/api/v1/' + casePortfolioType(type, false) + '/' + id;
    /* Send `DELETE` event with Ajax. */
    superagent
        .delete(target)
        .set('accept', 'json')
        .set('Authorization', jwt.token)
        .then((res)=>{
          if(res.errors){
            alert(res.errors);
            initPortfolio(currentUser, type);
          }
        })
        .catch(function (err) {
          alert(err);
          initPortfolio(currentUser, type);
        });
  });
}

/* Convert target label element into a input element while keeping the data. */
function loadPortfolioItem(index, item, type) {
  
  let form = document.getElementById(index);
  let label_id = '#' + item + '_label';
  let input_id = '#' + item + '_input';
  let label = form.querySelector(label_id);
  let input = form.querySelector(input_id);

  label.style.display = 'none';
  input.style.display = 'block';
  input.focus();

  if(item == 'category' || item == 'occurDate' || item == 'privacy'){
      input.addEventListener('change', ()=>{ updateAndSwitch(type);});
  }

  let updateBtn = form.querySelector('#update');
  updateBtn.addEventListener('click', (e)=>{ 
    e.preventDefault();
    updateAndSwitch(type);
  });

  input.addEventListener('keydown', (event)=> {
      /* Cancel input when pressing ESC. */
      if(event.keyCode == 27){
        event.preventDefault();
        label.style.display = 'block';
        input.style.display = 'none';
        input.value = label.innerText;
      }

      /* Update data when pressing ENTER. */
      if (event.keyCode === 13) {
        event.preventDefault();
        updateAndSwitch(type);
      }
  });

  function updateAndSwitch(type){
    isLogin();
    label.innerText = input.value;
    label.style.display = 'block';
    input.style.display = 'none';

    target = '/api/v1/' + casePortfolioType(type, false) + '/' + index;
    let data = {};
    if(item == 'occurDate'){
      let date = new Date(input.value);
      data[item] = date.toISOString();
    }
    else {
      data[item] = input.value;
    }
    /* Update the income item by sending a `PATCH` event with Ajax. */
    superagent
      .patch(target)
      .send(data)
      .set('accept', 'json')
      .set('Authorization', jwt.token)
      .then((res)=>{
        if(res.errors){
          alert(res.errors);
          initPortfolio(currentUser, type);
        }
      })
      .catch(function (err) {
        alert(err);
        initPortfolio(currentUser, type);
      });
  }
}

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

function addFriend(){
  isLogin();
  let tokenString = Session.get('token_type') + ' ' + Session.get('token');
  let url = window.location.href.split('/');
  let id = url[url.length-1];
  let addF = document.querySelector('#addFriend');
  let unF = document.querySelector('#unfriend');
  
  let mutation = 'mutation{\
      addFriend(id:"' + id + '")\
  }';
  superagent
  .post('/api/v1/graph')
  .set('accept', 'json')
  .set('Authorization', tokenString)
  .send({'query': mutation})
  .then(function (res) {
    addF.style.display = 'none';
    unF.style.display = 'block';
  })
  .catch(function (err) {
      alert(err);
  });
}

function unfriend(){
  isLogin();
  let tokenString = Session.get('token_type') + ' ' + Session.get('token');
  let url = window.location.href.split('/');
  let id = url[url.length-1];
  let addF = document.querySelector('#addFriend');
  let unF = document.querySelector('#unfriend');
  
  let mutation = 'mutation{\
    addFriend(id:"' + id + '")\
}';  
  superagent
  .post('/api/v1/graph')
  .set('accept', 'json')
  .set('Authorization', tokenString)
  .send({'query': mutation})
  .then(function (res) {
      addF.style.display = 'block';
      unF.style.display = 'none';
  })
  .catch(function (err) {
      alert(err);
      
  });
}