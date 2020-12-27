let baseURL = 'http://localhost:8080';

document.addEventListener('DOMContentLoaded', function () {
    /* Event listeners related to TODO creation. */
    (function () {
        /* The input entry listens to ENTER press event. */
        let form = document.querySelector('form');

        let input = form.querySelector('input');

        input.addEventListener('keydown', function (ev) {
            if (ev.which === 13) {
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

            superagent
                .post(`${baseURL}/todo/`)
                .send({
                    item: item,
                    index: 0
                })
                .set('accept', 'json')
                .then(function (res) {
                    clearMessage();

                    console.log(res.body);

                    addTODO(res.body);
                    input.value = '';
                })
                .catch(function (err) {
                    if (err.response) {
                        showMessage(err.response.message);
                    }
                });
        }
    })();

    /* The document listens to ESC press event */
    document.addEventListener('keydown', function (event) {
        if (event.which === 27) {
            let todos = document.getElementsByClassName('todo');

            for (let i = 0; i < todos.length; i++) {
                let label = todos[i].querySelector('label');
                let inputTODO = todos[i].querySelector('[name="index"]');
                let indexTODO = inputTODO.value;

                if (!label) {
                    let input = todos[i].querySelector('input');

                    let text = input.value;

                    let label = document.createElement('label');

                    label.classList.add('col-form-label');
                    label.innerText = text;

                    label.addEventListener('click', function () {
                        loadItem(indexTODO);
                    });

                    let index = todos[i].querySelector('[name="index"]').getAttribute('value');

                    let inputIndex = document.createElement('input');

                    inputIndex.setAttribute('value', index);
                    inputIndex.name = 'index'
                    inputIndex.setAttribute('hidden', true);

                    todos[i].innerHTML = '';
                    todos[i].appendChild(label);
                    todos[i].appendChild(inputIndex);
                }
            }
        }
    });

    /* Load initial TODO items */
    superagent
        .get(`${baseURL}/todos/`)
        .set('accept', 'json')
        .then(function (res) {
            clearMessage();

            let ts = res.body.todos;

            console.log(ts);

            for (let i = 0; i < ts.length; i++) {
                addTODO(ts[i]);
            }
        })
        .catch(function (err) {
            if (err.reponse) {
                showMessage(err.reponse.message);
            }
        });
});

/* Add a TODO item. */
function addTODO(todo) {
    let item = todo.item;
    let index = todo.index;

    /* Create a label element for item text. */
    let label = document.createElement('label');

    label.classList.add('col-form-label');
    label.innerText = item;

    label.addEventListener('click', function () {
        loadItem(index);
    });

    /* Create a hidden input element for item index. */
    let input = document.createElement('input');

    input.name = 'index';
    input.setAttribute('value', index);
    input.setAttribute('hidden', true);

    /* A holder for the label and input elements. */
    let row = document.createElement('div');

    row.classList.add('offset-lg-1');
    row.classList.add('col-lg-8');
    row.classList.add('offset-md-1');
    row.classList.add('col-md-7');
    row.classList.add('todo');

    row.style.marginTop = '5pt';
    row.style.marginBottom = '5pt';

    row.appendChild(label);
    row.appendChild(input);

    /* Create a `Update` button. */
    let btnUpdate = document.createElement('button');

    btnUpdate.innerText = 'Update';
    btnUpdate.type = 'submit';
    btnUpdate.name = '_method';
    btnUpdate.value = 'update';
    btnUpdate.addEventListener('click', function (ev) {
        ev.preventDefault();

        /* Get TODO item and index from the page. */
        let item;
        let index;

        let form = btnUpdate.parentNode.parentNode.parentNode;

        let todo = form.querySelector('.todo');

        let label = todo.querySelector('label');

        if (label) {
            item = label.innerText;
        } else {
            let _input = todo.querySelector('input');

            item = _input.value;
        }

        index = todo.querySelector('[name="index"]').getAttribute('value');

        console.log({
            item: item,
            index: Number(index)
        });

        /* Send `PUT` event with Ajax. */
        superagent
            .put(`${baseURL}/todo/`)
            .send({
                item: item,
                index: Number(index)
            })
            .set('accept', 'json')
            .then(function (res) {
                clearMessage();

                let form = btnUpdate.parentNode.parentNode.parentNode;

                let todo = form.querySelector('.todo');
                let inputTODO = todo.querySelector('[name="index"]');
                let indexTODO = inputTODO.value;

                let item = res.body.item;
                let index = res.body.index;

                /* Re-create new label and hidden input elements. */
                let _label = document.createElement('label');

                _label.classList.add('col-form-label');
                _label.innerText = item;

                _label.addEventListener('click', function () {
                    loadItem(indexTODO);
                });

                let inputIndex = document.createElement('input');

                inputIndex.setAttribute('value', index);
                inputIndex.name = 'index';
                inputIndex.setAttribute('hidden', true);

                /* Clear old elements and append new elements. */
                todo.innerHTML = '';
                todo.appendChild(_label);
                todo.appendChild(inputIndex);
            })
            .catch(function (err) {
                if (err.response) {
                    showMessage(err.response.message);
                }
            });
    }, false);

    btnUpdate.classList.add('btn');
    btnUpdate.classList.add('btn-secondary');

    /* Create a `Delete` button. */
    let btnDelete = document.createElement('button');

    btnDelete.innerText = 'Delete';
    btnDelete.type = 'submit';
    btnDelete.name = '_method';
    btnDelete.value = 'delete';
    btnDelete.addEventListener('click', function (ev) {
        ev.preventDefault();

        /* Get TODO index from the page. */
        let item;
        let index;

        let form = btnUpdate.parentNode.parentNode.parentNode;

        let todo = form.querySelector('.todo');

        let label = todo.querySelector('label');

        if (label) {
            item = label.innerText;
        } else {
            let _input = todo.querySelector('input');

            item = _input.value;
        }

        index = todo.querySelector('[name="index"]').getAttribute('value');

        console.log({
            item: item,
            index: Number(index)
        });

        /* Send `DELETE` event with Ajax. */
        superagent
            .delete(`${baseURL}/todo/`)
            .send({
                item: item,
                index: Number(index)
            })
            .set('accept', 'json')
            .then(function (res) {
                clearMessage();

                /* Remove the whole form. */
                let form = btnUpdate.parentNode.parentNode.parentNode;

                form.parentNode.removeChild(form);
            })
            .catch(function (err) {
                if (err.response) {
                    showMessage(err.response.message);
                }
            });
    }, false);

    btnDelete.classList.add('btn');
    btnDelete.classList.add('btn-info');

    /* A holder for `Update` and `Delete` buttons. */
    let rowButtons = document.createElement('div');

    rowButtons.classList.add('col-lg-3');
    rowButtons.classList.add('col-md-4');

    rowButtons.appendChild(btnUpdate);
    rowButtons.appendChild(btnDelete);

    /* Create a new HTML form. */
    let form = document.createElement('form');

    form.action = '/todo/';
    form.method = 'POST';

    let div = document.createElement('div');

    div.classList.add('row');

    div.appendChild(row);
    div.appendChild(rowButtons);

    form.appendChild(div);

    /* Append the created form to the index page. */
    let todoList = document.getElementById('todos');
    todoList.appendChild(form);
}

/* Convert target label element into a input element while keeping the data. */
function loadItem(index) {
    let todos = document.getElementsByClassName('todo');

    for (let i = 0; i < todos.length; i++) {
        let label = todos[i].querySelector('label');
        let inputTODO = todos[i].querySelector('[name="index"]');
        let indexTODO = inputTODO.value;

        /* Only convert the label element when the index is matched. */
        if (Number(indexTODO) === Number(index)) {
            /* Only convert it when `label` element exists. */
            if (label) {
                console.log('Convert label to input');

                let text = label.innerText;

                let input = document.createElement('input');

                input.classList.add('form-control');
                input.name = 'todo';
                input.setAttribute('value', text);

                input.addEventListener('keydown', function (event) {
                    /* Update data when pressing ENTER or ESC key. */
                    if (event.which === 13 || event.which === 27) {
                        let form = event.target.parentNode.parentNode.parentNode;

                        let todo = form.querySelector('.todo');

                        let _input = todo.querySelector('input');

                        let item = _input.value;
                        let index = todo.querySelector('[name="index"]').getAttribute('value');

                        console.log({
                            item: item,
                            index: Number(index)
                        });

                        /* Update the TODO item by sending a `PUT` event with Ajax. */
                        superagent
                            .put(`${baseURL}/todo/`)
                            .send({
                                item: item,
                                index: Number(index)
                            })
                            .set('accept', 'json')
                            .then(function (res) {
                                clearMessage();

                                let form = btnUpdate.parentNode.parentNode.parentNode;

                                let _todo = form.querySelector('.todo');
                                let _inputTODO = todo.querySelector('[name="index"]');
                                let _indexTODO = _inputTODO.value;

                                let item = res.body.item;
                                let index = res.body.index;

                                /* Re-create new label and input elements. */
                                let _label = document.createElement('label');

                                _label.classList.add('col-form-label');
                                _label.innerText = item;

                                _label.addEventListener('click', function () {
                                    loadItem(_indexTODO);
                                });

                                let inputIndex = document.createElement('input');

                                inputIndex.setAttribute('value', index);
                                inputIndex.name = 'index'
                                inputIndex.setAttribute('hidden', true);

                                /* Clear old elements and append new elements. */
                                _todo.innerHTML = '';
                                _todo.appendChild(_label);
                                _todo.appendChild(inputIndex);
                            })
                            .catch(function (err) {
                                if (err.response) {
                                    showMessage(err.response.message);
                                }
                            });
                    }
                });

                let index = todos[i].querySelector('[name="index"]').getAttribute('value');

                console.log(todos[i].querySelector('[name="index"]'));
                console.log(`index: ${index}`);

                let inputIndex = document.createElement('input');

                inputIndex.setAttribute('value', index);
                inputIndex.name = 'index'
                inputIndex.setAttribute('hidden', true);

                todos[i].innerHTML = '';
                todos[i].appendChild(input);
                todos[i].appendChild(inputIndex);
            }
        } else {
            /* Convert input elments back to label elements while the element is not the target. */
            if (!label) {
                let input = todos[i].querySelector('input');

                let text = input.getAttribute('value');

                /* Re-create label and input elements. */
                let label = document.createElement('label');
                let inputTODO = todos[i].querySelector('[name="index"]');
                let indexTODO = inputTODO.value;

                label.classList.add('col-form-label');
                label.innerText = text;

                label.addEventListener('click', function () {
                    loadItem(indexTODO);
                });

                let index = todos[i].querySelector('[name="index"]').getAttribute('value');

                let inputIndex = document.createElement('input');

                inputIndex.setAttribute('value', index);
                inputIndex.name = 'index'
                inputIndex.setAttribute('hidden', true);

                /* Clear old elements and append new elements. */
                todos[i].innerHTML = '';
                todos[i].appendChild(label);
                todos[i].appendChild(inputIndex);
            }
        }
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
