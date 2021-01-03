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

 async function checkMyFriend(){
  checkToken(jwt);
  let tokenString = Session.get('token_type') + ' ' + Session.get('token');
  let url = window.location.href.split('/');
  let id = url[url.length-1];

  query = '{\
      myFriends{\
        id\
      }\
    }';

  await superagent
  .post('/api/v1/graph')
  .set('accept', 'json')
  .set('Authorization', tokenString)
  .send({'Query': query,})
  .then(function (res) {
      clearMessage()
      let out = res.body.data.MyFriends;
      for(let i = 0; i < out.length; i++){
        if(out[i].Id == id){
          document.querySelector('#addFriend').display = 'none';
          document.querySelector('#unfriend').display = 'block';
        }
      }
      input.value = '';
  })
  .catch(function (err) {
      if (err.response) {
          showMessage(err.response.message);
      }
  });
 }

//get UserId through gql API
async function getUserProfile(currentUser){
  checkToken(jwt);
  let tokenString = Session.get('token_type') + ' ' + Session.get('token');
  let url = window.location.href.split('/');
  let id = url[url.length-1];

  query = '{\
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
      }\
    }';
  console.log(query);

  await superagent
  .post('/api/v1/graph')
  .set('accept', 'json')
  .set('Authorization', tokenString)
  .send({'Query': query,})
  .then(function (res) {
      clearMessage()
      let out = res.body.data.getUser;
      currentUser.id = out.id;
      currentUser.userId = out.userId;
      currentUser.nickName = out.nickName;
      currentUser.email = out.email;
      currentUser.friends = [...out.friends];
      currentUser.followers = [...out.followers];

      let resp = document.getElementById('response');
      resp.innerHTML = JSON.stringify(res.body);
      console.log(res.body);

      input.value = '';
  })
  .catch(function (err) {
      if (err.response) {
          showMessage(err.response.message);
      }
  });
}

async function getAccessToken(jwt){
await superagent
  .get('/api/v1/auth/ics/refresh_token')
  .set('accept', 'json')
  .then(function (res) {
      clearMessage();
      console.log(res.body);
      Session.set('token', res.body.token);
      Session.set('token_expiry', res.body.token_expiry);
      Session.set('token_type', res.body.type);
      jwt.token = res.body.token;
      jwt.token_expiry = res.body.token_expiry;
      jwt.token_type = res.body.type;
      input.value = '';
  })
  .catch(function (err) {
    
    if (err.response) {
      console.log(err);
      alert(err.response.message);
      logout();
    }
    
  });
}

function checkToken(jwt){
  jwt.token = Session.get('token');
  jwt.token_expiry = new Date(Session.get('token_expiry'));
  jwt.token_type = Session.get('token_type');
  var now = new Date()
  
  if (typeof(jwt.token) == 'undefined' || jwt.token_expiry < now){
    getAccessToken(jwt);
  }
  if(jwt.token_expiry > now){
    document.querySelector('#login').style.display = 'none';
    document.querySelector('#signup').style.display = 'none';
    document.querySelector('#logout').style.display = 'block';
  }
  
}

function logout(){
  console.log('logout')
  Session.set('token_expiry', Date());
  jwt.token_expiry = Date();

  //ask server to set refresh token invalid
  var request = new XMLHttpRequest(); 
  request.open('GET', '/api/v1/auth/ics/logout', true);
  request.send();
  request.onerror = showMessage();

  document.querySelector('#login').style.display = 'block';
  document.querySelector('#signup').style.display = 'block';
  document.querySelector('#logout').style.display = 'none';
}

var jwt = {};
checkToken(jwt);


var currentUser = {};
getUserProfile(currentUser).then(()=>{
  initPortfolio(currentUser, INCOME);
  initPortfolio(currentUser, COST);
});

const COST = 'COST';
const INCOME = 'INCOME';

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
  <div class='col-md-7 border'>\
  <label for='description' id='description_label' style='width:100%;'>" + casePortfolioType(type, true) + " Description</label>\
  <input type='text' class='form-control' name='description' placeholder='" + casePortfolioType(type, true) + " Description' id='description_input' style='display: none;'>\
  </div>\
  <div class='col-md-3 border'>\
  <label for='amount' id='amount_label' style='width:100%;'>Amount: $888</label>\
  <input type='text' class='form-control' name='amount' placeholder='888' id='amount_input' style='display: none;'>\
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

document.addEventListener('DOMContentLoaded', function () {
  /* Event listeners related to portfolio creation. */
  (function () {
      /* The input entry listens to ENTER press event. */
      let form = document.querySelector('#add');
      let description = document.querySelector('#description');
      let amount = document.querySelector('#amount');
      let incomeCategory = document.querySelector('#incomeCat');
      let costCategory = document.querySelector('#costCat');
      let occurDate = document.querySelector('#occurDate')
      let cost = document.querySelector('#cost');
      let income = document.querySelector('#income');
      let category;
      
      description.addEventListener('keydown', function (ev) {
          if (ev.keyCode === 13) {
            ev.preventDefault();
            if(cost.checked){
                category = costCategory;
                createPortfolio(COST);
            }
            else if(income.checked){
                category = incomeCategory;
                createPortfolio(INCOME);
            }
          }
      }, false);

      amount.addEventListener('keydown', function (ev) {
        if (ev.keyCode === 13) {
            ev.preventDefault();
            if(cost.checked){
                category = costCategory;
                createPortfolio(COST);
            }
            else if(income.checked){
                category = incomeCategory;
                createPortfolio(INCOME);
            }
        }
      }, false);

      /* The `Add` button listens to click event. */
      let btn = form.querySelector('#create');

      btn.addEventListener('click', function (ev) {
          ev.preventDefault();
          if(cost.checked){
              category = costCategory;
              createPortfolio(COST);
          }
          else if(income.checked){
              category = incomeCategory;
              createPortfolio(INCOME);
          }
      }, false);

      async function createPortfolio(type) {
        checkToken(jwt);
        let target = '/api/v1/' + casePortfolioType(type, false);
        let data = {};
        data.description = description.value;
        data.amount = amount.value;
        data.category = category.value;
        let date = new Date(occurDate.value);
        data.occurDate = date.toISOString();
            
        await superagent
            .post(target)
            .send(data)
            .set('accept', 'json')
            .set('Authorization', jwt.token)
            .then(function (res) {
                
                clearMessage();
                if(type == COST){
                  addPortfolio(type, res.body.CreateCost);
                }
                else{
                  addPortfolio(type, res.body.CreateIncome);
                }
            })
            .catch(function (err) {
                if (err.response) {
                    showMessage(err.response.message);
                }
            });
    }
  })();
})

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
          clearMessage();
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
          if (err.reponse) {
              showMessage(err.reponse.message);
          }
      });
}

/* Add a Portfolio item. */
function addPortfolio(res, type) {
  let id = res.Id;
  let description = res.Description;
  let amount = res.Amount;
  let category = res.Category;
  let date = new Date(res.OccurDate);
  let date_month = date.getMonth()+1;
  let occurDate = date.getFullYear() + '-' + date_month + '-' + date.getDate();
  let vote = res.Vote.length;

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
    checkToken(jwt);
    e.preventDefault();
    let target = '/api/v1/' + casePortfolioType(type, false) + '/vote/' + id;
    superagent
        .put(target)
        .set('accept', 'json')
        .set('Authorization', jwt.token)
        .then(function (res) {
            clearMessage();
            input = '';
            vote_label.innerText = (type==INCOME? res.body.VoteIncome:res.body.VoteCost);
        })
        .catch(function (err) {
            if (err.response) {
                showMessage(err.response.message);
                initPortfolio(currentUser, type);
            }
        });
  })

  let deleteBtn = form.querySelector('#delete');
  deleteBtn.addEventListener('click', (e)=>{
    checkToken(jwt);
    e.preventDefault();
    form.parentNode.removeChild(form);
    let target = '/api/v1/' + casePortfolioType(type, false) + '/' + id;
    /* Send `DELETE` event with Ajax. */
    superagent
        .delete(target)
        .set('accept', 'json')
        .set('Authorization', jwt.token)
        .then(function () {
            clearMessage();
        })
        .catch(function (err) {
            if (err.response) {
                showMessage(err.response.message);
                initPortfolio(currentUser, type);
            }
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

  if(item == 'category' || item == 'occurDate'){
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
    checkToken(jwt);
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
        .then(function () {
            clearMessage();
        })
        .catch(function (err) {
            if (err.response) {
                showMessage(err.response.message);
                initPortfolio(currentUser, type);
            }
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



