let token = 'token';
getAccessToken();
document.addEventListener('DOMContentLoaded', function () {
    /* Event listeners related to TODO creation. */
    (function () {
        /* The input entry listens to ENTER press event. */
        let form = document.querySelector('#add');

        let input = form.querySelector('input');

        input.addEventListener('keydown', function (ev) {
            if (ev.key === 13) {
                ev.preventDefault();
                createTODO();
            }
        }, false);

        /* The `Add` button listens to click event. */
        let btn = form.querySelector('button');

        btn.addEventListener('click', function (ev) {
            ev.preventDefault();
            createTODO();
        }, false);
        function createTODO() {
            let item = input.value;
            let target = `/api/v1/user/` + document.querySelector('#id').value
        
            superagent
                .patch(target)
                .send({
                    item: item,
                    index: 0
                })
                .set('accept', 'json')
                .set('Authorization', token)
                .then(function (res) {
                    clearMessage();
        
                    console.log(res.body);
                    input.value = '';
                })
                .catch(function (err) {
                    if (err.response) {
                        showMessage(err.response.message);
                    }
                });
        }
        
    })();
})


function getAccessToken(){
    const res = superagent
        .get('/api/v1/auth/ics/refresh_token')
        .set('accept', 'json')
        .then(function (res) {
            clearMessage();

            token = res.body.token
            console.log(token);

            input.value = '';
        })
        .catch(function (err) {
            if (err.response) {
                showMessage(err.response.message);
            }
        });
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