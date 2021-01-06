import {StringValidator} from './Validation';
import {LettersOnlyValidator} from './LettersOnlyValidator';
import {ZipCodevalidator} from './ZipCodevalidator';

let strings = ["Hello","98052","101"];


let validators : {[s:string]: StringValidator; } = {};

validators["zip code"] = new ZipCodevalidator();

validators["Letters only"] = new LettersOnlyValidator();

strings.forEach(s =>{
    for(let name in validators){
        console.log("s:"+s+"name:"+name+",match:"+validators[name].isAcceptable(s));
    }
});