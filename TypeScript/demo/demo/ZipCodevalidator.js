"use strict";
exports.__esModule = true;
exports.ZipCodevalidator = void 0;
var numberRegex = /^[0-9]+$/;
var ZipCodevalidator = /** @class */ (function () {
    function ZipCodevalidator() {
    }
    ZipCodevalidator.prototype.isAcceptable = function (s) {
        return s.length === 5 && numberRegex.test(s);
    };
    return ZipCodevalidator;
}());
exports.ZipCodevalidator = ZipCodevalidator;
