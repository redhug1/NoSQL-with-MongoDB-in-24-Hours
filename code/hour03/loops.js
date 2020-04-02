var weekDays = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday"];
print("Week Days: ");
printjson(weekDays)
print("Combining Arrays: ");
var fullWeek = weekDays.slice(0);
fullWeek.unshift("Sunday");
fullWeek.push("Saturday");
print("Full Week: ");
printjson(fullWeek);
var midWeek = weekDays.splice(1, 3);
print("Mid Week: ");
printjson(midWeek);
print("Sliced weekdays: ");
printjson(weekDays);