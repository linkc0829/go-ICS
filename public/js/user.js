function addFriend(){
    checkToken(jwt);
    let tokenString = Session.get('token_type') + ' ' + Session.get('token');
    let url = window.location.href.split('/');
    let id = url[url.length-1];
    let addFriend = document.querySelector('#addFriend');
    let unfriend = document.querySelector('#unfriend');
    addFriend.style.display = 'none';
    unfriend.style.display = 'block';
    mutation = '{\
        addFriend(id: ' + id + '){\
            id\
        }\
    }';    
    superagent
    .post('/api/v1/graph')
    .set('accept', 'json')
    .set('Authorization', tokenString)
    .send({'Mutation': mutation,})
    .then(function () {
        clearMessage()
        input.value = '';
    })
    .catch(function (err) {
        if (err.response) {
            showMessage(err.response.message);
        }
        addFriend.style.display = 'block';
        unfriend.style.display = 'none';
    });
}

function unfriend(){
    checkToken(jwt);
    let tokenString = Session.get('token_type') + ' ' + Session.get('token');
    let url = window.location.href.split('/');
    let id = url[url.length-1];
    let addFriend = document.querySelector('#addFriend');
    let unfriend = document.querySelector('#unfriend');
    addFriend.style.display = 'block';
    unfriend.style.display = 'none';
    mutation = '{\
        addFriend(id: ' + id + '){\
            id\
        }\
    }';    
    superagent
    .post('/api/v1/graph')
    .set('accept', 'json')
    .set('Authorization', tokenString)
    .send({'Mutation': mutation,})
    .then(function () {
        clearMessage()
        input.value = '';
    })
    .catch(function (err) {
        if (err.response) {
            showMessage(err.response.message);
        }
        addFriend.style.display = 'none';
        unfriend.style.display = 'block';
    });
}