function formatGreeting(name, city){
    var retStr = "";
    retStr += "Hello " + name + "\n";
    retStr += "Welcome to " + city + "!";
    return retStr
}
var greeting = formatGreeting("Frodo", "Rivendell");
print(greeting);
greeting = formatGreeting("Arthur", "Camelot");
print(greeting);