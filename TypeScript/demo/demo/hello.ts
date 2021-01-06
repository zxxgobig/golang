
function greeter<T>(person: T): string {
    let str = "Hello,"+person;
    return str
}

let user = "Jane User ";
document.body.innerHTML = greeter<string>(user);

function getProperty(obj: any, key: any) {
    return obj[key];
}

let x = { a: 1, b: 2, c: 3, d: 4 };

getProperty(x, "a"); // okay
getProperty(x, "m"); // error: Argument of type 'm' isn't assignable to 'a' | 'b' | 'c' | 'd'.
