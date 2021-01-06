"use strict";
exports.__esModule = true;
var LettersOnlyValidator_1 = require("./LettersOnlyValidator");
var ZipCodevalidator_1 = require("./ZipCodevalidator");
var strings = ["Hello", "98052", "101"];
var validators = {};
validators["zip code"] = new ZipCodevalidator_1.ZipCodevalidator();
validators["Letters only"] = new LettersOnlyValidator_1.LettersOnlyValidator();
strings.forEach(function (s) {
    for (var name_1 in validators) {
        console.log("s:" + s + "name:" + name_1 + ",match:" + validators[name_1].isAcceptable(s));
    }
});
