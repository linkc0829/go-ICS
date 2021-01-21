db = db.getSiblingDB('admin')
db.auth(user, pwd)

//https://gist.github.com/solenoid/1372386
// var mongoObjectId = function () {
//     var timestamp = (new Date().getTime() / 1000 | 0).toString(16);
//     return timestamp + 'xxxxxxxxxxxxxxxx'.replace(/[x]/g, function() {
//         return (Math.random() * 16 | 0).toString(16);
//     }).toLowerCase();
// };

db = db.getSiblingDB('ics')
db.users.insert({
    userid: "admin",
    email: "admin@icsharing.com",
    nickname: "admin",
    createAt: new Date(),
    friends: [],
    lastIncomeQuery: new Date(),
    lastCostQuery: new Date(),
    provider: "ics",
    role: "ADMIN"
})