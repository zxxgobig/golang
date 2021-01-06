import {StringValidator} from './Validation';

const lettersRegex = /^[A-Za-z]+$/;

export class LettersOnlyValidator implements StringValidator {
    isAcceptable(s:string):boolean {
        return lettersRegex.test(s);
    }

}