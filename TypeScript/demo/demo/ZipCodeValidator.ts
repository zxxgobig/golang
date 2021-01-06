import {StringValidator} from './Validation';

const numberRegex = /^[0-9]+$/;

export class ZipCodevalidator implements StringValidator{
    isAcceptable(s:string){
        return s.length ===5 && numberRegex.test(s);
    }
}