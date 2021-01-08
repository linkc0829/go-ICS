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
                  createPortfolio(COST).then((res)=>{
                    initPortfolio(currentUser, COST);
                    clearForm();
                  })
              }
              else if(income.checked){
                  category = incomeCategory;
                  createPortfolio(INCOME).then((res)=>{
                    initPortfolio(currentUser, INCOME);
                    clearForm();
                  })
              }
            }
        }, false);
  
        amount.addEventListener('keydown', function (ev) {
          if (ev.keyCode === 13) {
              ev.preventDefault();
              if(cost.checked){
                  category = costCategory;
                  createPortfolio(COST).then((res)=>{
                    initPortfolio(currentUser, COST);
                    clearForm();
                  })
              }
              else if(income.checked){
                  category = incomeCategory;
                  createPortfolio(INCOME).then((res)=>{
                    initPortfolio(currentUser, INCOME);
                    clearForm();
                  })
              }
          }
        }, false);
  
        /* The `Add` button listens to click event. */
        let btn = form.querySelector('#create');
  
        btn.addEventListener('click', function (ev) {
            ev.preventDefault();
            if(cost.checked){
                category = costCategory;
                createPortfolio(COST).then((res)=>{
                  initPortfolio(currentUser, COST);
                  clearForm();
                })
                
            }
            else if(income.checked){
                category = incomeCategory;
                createPortfolio(INCOME).then((res)=>{
                  initPortfolio(currentUser, INCOME);
                  clearForm();
                })
            }
        }, false);
  
        async function createPortfolio(type) {
          isLogin();
          let target = '/api/v1/' + casePortfolioType(type, false);
          let data = {};
          data.description = description.value;
          data.amount = amount.value;
          data.category = category.value;
          let date = new Date(occurDate.value);
          data.occurDate = date.toISOString();
  
          if(data.description == "" || amount == "" ||  data.category == "ZERO"){
            alert('input incomplete');
            return;
          }
          if(isNaN(amount)){
            alert('amount must be numbers');
            return;
          }
          if(date < new Date()){
            alert('OccurDate must in the future.')
            return;
          }
              
          await superagent
            .post(target)
            .send(data)
            .set('accept', 'json')
            .set('Authorization', jwt.token)
            .then(function (res) {            
                if(type == COST){
                  addPortfolio(type, res.body.CreateCost);
                }
                else{
                  addPortfolio(type, res.body.CreateIncome);
                }
            })
            .catch(function (err) {
              alert(err);
            });
        }
        function clearForm(){
          document.querySelector('#description').value = "";
          document.querySelector('#amount').value = "";
          document.querySelector('#incomeCat').value = "ZERO";
          document.querySelector('#costCat').value = "ZERO";
          document.querySelector('#occurDate').value = "";
        }
    })();
  })