"use strict";
exports.__esModule = true;
exports.LettersOnlyValidator = void 0;
var lettersRegex = /^[A-Za-z]+$/;
var LettersOnlyValidator = /** @class */ (function () {
    function LettersOnlyValidator() {
    }
    LettersOnlyValidator.prototype.isAcceptable = function (s) {
        return lettersRegex.test(s);
    };
    return LettersOnlyValidator;
}());
exports.LettersOnlyValidator = LettersOnlyValidator;
