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

//get UserId through gql API
async function getMyProfile(me){
  let now = new Date();
  let token_expiry = Date.parse(Session.get('token_expiry'));
  if(now >= token_expiry){
      getAccessToken();
  }

  let tokenString = Session.get('token_type') + ' ' + Session.get('token');
  query = '{\
      me{\
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

  await superagent
  .post('/api/v1/graph')
  .set('accept', 'json')
  .set('Authorization', tokenString)
  .send({'Query': query,})
  .then(function (res) {
      clearMessage()
      let out = res.body.data.me;
      me.id = out.id;
      me.userId = out.userId;
      me.nickName = out.nickName;
      me.email = out.email;
      me.friends = [...out.friends];
      me.followers = [...out.followers];

      let resp = document.getElementById('response');
      resp.innerHTML = JSON.stringify(res.body);

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
          showMessage(err.response.message);
      }
  });
}

var jwt = {};
jwt.token = Session.get('token');
jwt.token_expiry = new Date(Session.get('token_expiry'));
jwt.token_type = Session.get('token_type');
var now = new Date()

if (typeof(jwt.token) == 'undefined' || jwt.token_expiry < now){
getAccessToken(jwt);
}

var me = {};
getMyProfile(me).then(()=>{
  initIncome(me);
  initCost(me);
});


const incomeFormInnerHtml = "<div class='form-row'>\
<div class='col-md-7 border'>\
<label for='description' id='description_label' style='width:100%;'>Income Description</label>\
<input type='text' class='form-control' name='description' placeholder='Income Description' id='description_input' style='display: none;'>\
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
<div class='col-md-3 border'>\
<label for='category' id='category_label' style='width:100%;'>INVESTMENT</label>\
<select name='category' id='category_input' class='custom-select' style='display: none;'>\
  <option selected>Choose a category</option>\
  <option value='INVESTMENT'>INVESTMENT</option>\
  <option value='SALARY'>SALARY</option>\
  <option value='PARTTIME'>PART TIME</option>\
  <option value='OTHERS'>OTHERS</option>\
</select>\
</div>\
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

const costFormInnerHtml = "<div class='form-row'>\
<div class='col-md-7 border'>\
<label for='description' id='description_label' style='width:100%;'>Cost Description</label>\
<input type='text' class='form-control' name='description' placeholder='Income Description' id='description_input' style='display: none;'>\
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
<div class='col-md-3 border'>\
<label for='category' id='category_label' style='width:100%;'>LEARNING</label>\
<select name='category' id='costCat' style='display: none;' class='custom-select'>\
  <option selected>Choose a category</option>\
  <option value='INVESTMENT'>INVESTMENT</option>\
  <option value='DAILY'>DAILY</option>\
  <option value='LEARNING'>LEARNING</option>\
  <option value='CHARITY'>CHARITY</option>\
  <option value='OTHERS'>OTHERS</option>\
</select>\
</div>\
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
                  createCost();
              }
              else if(income.checked){
                  category = incomeCategory;
                  createIncome();
              }
          }
      }, false);

      amount.addEventListener('keydown', function (ev) {
          if (ev.keyCode === 13) {
              ev.preventDefault();
              if(cost.checked){
                  category = costCategory;
                  createCost();
              }
              else if(income.checked){
                  category = incomeCategory;
                  createIncome();
              }
          }
      }, false);

      /* The `Add` button listens to click event. */
      let btn = form.querySelector('#create');

      btn.addEventListener('click', function (ev) {
          ev.preventDefault();
          if(cost.checked){
              category = costCategory;
              createCost();
          }
          else if(income.checked){
              category = incomeCategory;
              createIncome();
          }
      }, false);

      async function createIncome() {
          let target = `/api/v1/income`;
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
                  addIncome(res.body)
              })
              .catch(function (err) {
                  if (err.response) {
                      showMessage(err.response.message);
                  }
              });
      }
      async function createCost() {
        let target = '/api/v1/cost';
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
                addCost(res.body)
            })
            .catch(function (err) {
                if (err.response) {
                    showMessage(err.response.message);
                }
            });
    }
  })();
})

//initial incomes or reload all records
async function initIncome(me){
  let target = '/api/v1/user/' + me.id + '/income';
  await superagent
      .get(target)
      .set('accept', 'json')
      .set('Authorization', jwt.token)
      .then(function (res) {
          clearMessage();
          let incomes = res.body.GetUserIncome;

          for (let i = 0; i < incomes.length; i++) {
              addIncome(incomes[i]);
          }
      })
      .catch(function (err) {
          if (err.reponse) {
              showMessage(err.reponse.message);
          }
      });
}

//initial icosts or reload all records
async function initCost(me){
  let target = '/api/v1/user/' + me.id + '/cost';
  await superagent
      .get(target)
      .set('accept', 'json')
      .set('Authorization', jwt.token)
      .then(function (res) {
          clearMessage();
          let costs = res.body.GetUserCost;

          for (let i = 0; i < costs.length; i++) {
              addCost(costs[i]);
          }
      })
      .catch(function (err) {
          if (err.reponse) {
              showMessage(err.reponse.message);
          }
      });
}

/* Add a Income item. */
function addIncome(res) {
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
  form.innerHTML = incomeFormInnerHtml;
  let incomeList = document.querySelector('#income-list');
  incomeList.appendChild(form);

  let description_label = form.querySelector('#description_label');
  description_label.innerText = description;
  description_label.htmlFor = 'description_' + id;
  description_label.addEventListener('click', ()=>{ loadIncomeItem(id, 'description')})
  let description_input = form.querySelector('#description_input');
  description_input.name = 'description_' + id;
  description_input.value = description;

  let amount_label = form.querySelector('#amount_label');
  amount_label.innerText = amount;
  amount_label.htmlFor = 'amount_' + id;
  amount_label.addEventListener('click', ()=>{ loadIncomeItem(id, 'amount')})
  let amount_input = form.querySelector('#amount_input');
  amount_input.name = 'amount_' + id;
  amount_input.value = amount;

  let category_label = form.querySelector('#category_label');
  category_label.innerText = category;
  category_label.htmlFor = 'category_' + id;
  category_label.addEventListener('click', ()=>{ loadIncomeItem(id, 'category')})
  let category_input = form.querySelector('#category_input');
  category_input.name = 'category_' + id;
  category_input.value = category;

  let occurDate_label = form.querySelector('#occurDate_label');
  occurDate_label.innerText = occurDate;
  occurDate_label.htmlFor = 'occueDate_' + id;
  occurDate_label.addEventListener('click', ()=>{ loadIncomeItem(id, 'occurDate')})
  let occurDate_input = form.querySelector('#occurDate_input');
  occurDate_input.name = 'occueDate_' + id;
  occurDate_input.value = occurDate;

  let vote_label = form.querySelector('#vote_label');
  vote_label.innerText = vote;
  let voteBtn = form.querySelector('#vote');
  voteBtn.addEventListener('click', (e)=>{
      e.preventDefault();
      //++vote;
      //vote_label.innerText = vote;
      let target = '/api/v1/income/vote/' + id;
      superagent
          .put(target)
          .set('accept', 'json')
          .set('Authorization', jwt.token)
          .then(function (res) {
              clearMessage();
              input = '';
              vote_label.innerText = res.body.VoteIncome;
          })
          .catch(function (err) {
              if (err.response) {
                  showMessage(err.response.message);
                  //reloadAll()
              }
          });
  })

  let deleteBtn = form.querySelector('#delete');
  deleteBtn.addEventListener('click', (e)=>{
      e.preventDefault();
      form.parentNode.removeChild(form);
      let target = '/api/v1/income/' + id;
      /* Send `DELETE` event with Ajax. */
      superagent
          .delete(target)
          .set('accept', 'json')
          .set('Authorization', jwt.token)
          .then(function (res) {
              clearMessage();
          })
          .catch(function (err) {
              if (err.response) {
                  showMessage(err.response.message);
                  //reloadAll()
              }
          });
  });
}

/* Convert target label element into a input element while keeping the data. */
function loadIncomeItem(index, item) {
  
  let form = document.getElementById(index);
  let label_id = '#' + item + '_label';
  let input_id = '#' + item + '_input';
  let label = form.querySelector(label_id);
  let input = form.querySelector(input_id);

  label.style.display = 'none';
  input.style.display = 'block';
  input.focus();

  if(item == 'category' || item == 'occurDate'){
      input.addEventListener('change', ()=>{ updateAndSwitch();});
  }

  let updateBtn = form.querySelector('#update');
  updateBtn.addEventListener('click', (e)=>{ 
    e.preventDefault();
    updateAndSwitch();
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
        updateAndSwitch();
      }
  });

  function updateAndSwitch(){
      label.innerText = input.value;
      label.style.display = 'block';
      input.style.display = 'none';

      target = '/api/v1/income/' + index;
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
          .then(function (res) {
              clearMessage();
          })
          .catch(function (err) {
              if (err.response) {
                  showMessage(err.response.message);
                  //reloadAll()
              }
          });
  }
}

/* Add a Cost item. */
function addCost(res) {
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
  form.innerHTML = costFormInnerHtml;
  let costList = document.querySelector('#cost-list');
  costList.appendChild(form);

  let description_label = form.querySelector('#description_label');
  description_label.innerText = description;
  description_label.htmlFor = 'description_' + id;
  description_label.addEventListener('click', ()=>{ loadCostItem(id, 'description')})
  let description_input = form.querySelector('#description_input');
  description_input.name = 'description_' + id;
  description_input.value = description;

  let amount_label = form.querySelector('#amount_label');
  amount_label.innerText = amount;
  amount_label.htmlFor = 'amount_' + id;
  amount_label.addEventListener('click', ()=>{ loadCostItem(id, 'amount')})
  let amount_input = form.querySelector('#amount_input');
  amount_input.name = 'amount_' + id;
  amount_input.value = amount;

  let category_label = form.querySelector('#category_label');
  category_label.innerText = category;
  category_label.htmlFor = 'category_' + id;
  category_label.addEventListener('click', ()=>{ loadCostItem(id, 'category')})
  let category_input = form.querySelector('#category_input');
  category_input.name = 'category_' + id;
  category_input.value = category;

  let occurDate_label = form.querySelector('#occurDate_label');
  occurDate_label.innerText = occurDate;
  occurDate_label.htmlFor = 'occueDate_' + id;
  occurDate_label.addEventListener('click', ()=>{ loadCostItem(id, 'occurDate')})
  let occurDate_input = form.querySelector('#occurDate_input');
  occurDate_input.name = 'occueDate_' + id;
  occurDate_input.value = occurDate;

  let vote_label = form.querySelector('#vote_label');
  vote_label.innerText = vote;
  let voteBtn = form.querySelector('#vote');
  voteBtn.addEventListener('click', (e)=>{
      e.preventDefault();
      //++vote;
      //vote_label.innerText = vote;
      let target = '/api/v1/cost/vote/' + id;
      superagent
          .put(target)
          .set('accept', 'json')
          .set('Authorization', jwt.token)
          .then(function (res) {
              clearMessage();
              input = '';
              vote_label.innerText = res.body.VoteCost;
          })
          .catch(function (err) {
              if (err.response) {
                  showMessage(err.response.message);
                  //reloadAll()
              }
          });
  })

  let deleteBtn = form.querySelector('#delete');
  deleteBtn.addEventListener('click', (e)=>{
      e.preventDefault();
      form.parentNode.removeChild(form);
      let target = '/api/v1/cost/' + id;
      /* Send `DELETE` event with Ajax. */
      superagent
          .delete(target)
          .set('accept', 'json')
          .set('Authorization', jwt.token)
          .then(function (res) {
              clearMessage();
          })
          .catch(function (err) {
              if (err.response) {
                  showMessage(err.response.message);
                  //reloadAll()
              }
          });
  });
}

/* Convert target label element into a input element while keeping the data. */
function loadCostItem(index, item) {
  
  let form = document.getElementById(index);
  let label_id = '#' + item + '_label';
  let input_id = '#' + item + '_input';
  let label = form.querySelector(label_id);
  let input = form.querySelector(input_id);

  label.style.display = 'none';
  input.style.display = 'block';
  input.focus();

  if(item == 'category' || item == 'occurDate'){
      input.addEventListener('change', ()=>{ updateAndSwitch();});
  }

  let updateBtn = form.querySelector('#update');
  updateBtn.addEventListener('click', (e)=>{ 
    e.preventDefault();
    updateAndSwitch();
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
        updateAndSwitch();
      }
  });

  function updateAndSwitch(){
      label.innerText = input.value;
      label.style.display = 'block';
      input.style.display = 'none';

      target = '/api/v1/cost/' + index;
      let data = {};
      if(item == 'occurDate'){
        let date = new Date(input.value);
        data[item] = date.toISOString();
      }
      else {
        data[item] = input.value;
      }

      /* Update the cost item by sending a `PATCH` event with Ajax. */
      superagent
          .patch(target)
          .send(data)
          .set('accept', 'json')
          .set('Authorization', jwt.token)
          .then(function (res) {
              clearMessage();
          })
          .catch(function (err) {
              if (err.response) {
                  showMessage(err.response.message);
                  //reloadAll()
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



