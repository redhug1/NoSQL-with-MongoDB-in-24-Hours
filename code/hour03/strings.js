var myStr = "I think therefore I am.";
print("Original string: ");
print(myStr);
print("Finding the substring thing: ");
if (myStr.indexOf("think") != -1){
    print(myStr + " contains think");
}
print("Replacing the substring think with feel: ");
var newStr = myStr.replace("think", "feel");
print(newStr);
print("Converting the phrase into an array: ");
var myArr = myStr.split(" ");
printjson(myArr);