var vowelArr = "aeiou";
var consonantArr = "bcdfghjklmnpqrstvwxyz";
//var words = "the,be,and,of,to,it,I,can't,shouldn't,say,middle-class,apology,till";
//var wordArr = words.split(",");
// load in file that came from : https://github.com/first20hours/google-10000-english
var file = cat('google-10000-english.txt'); // read the file
var wordArr = file.split('\n'); // create array of words
var wordObjArr = new Array();
for (var i=0; i<wordArr.length; i++){
    try{
        var word = wordArr[i].toLowerCase();
        var vowelCnt = ("|"+word+"|").split(/[aeiou]/i).length-1;
        var consonantCnt = ("|"+word+"|").split(/[bcdfghjklmnpqrstvwxyz]/i).length-1;
        var letters = [];
        var vowels = [];
        var consonants = [];
        var other = [];
        for (var j=0; j<word.length; j++){
            var ch = word[j];
            if (letters.indexOf(ch) === -1){
                letters.push(ch);
            }
            if (vowelArr.indexOf(ch) !== -1){
                if (vowels.indexOf(ch) === -1){
                    vowels.push(ch);
                }
            } else if (consonantArr.indexOf(ch) !== -1){
                if (consonants.indexOf(ch) === -1){
                    consonants.push(ch);
                }
            } else {
                if (other.indexOf(ch) === -1){
                    other.push(ch);
                }
            }
        }
        var charsets = [];
        if (consonants.length){
            charsets.push({type:"consonants", chars:consonants});
        }
        if (vowels.length){
            charsets.push({type:"vowels", chars:vowels});
        }
        if (other.length){
            charsets.push({type:"other", chars:other});
        }
        var wordObj = {
            word: word,
            first: word[0],
            last: word[word.length-1],
            size: word.length,
            letters: letters,
            stats: {vowels: vowelCnt, consonants: consonantCnt},
            charsets: charsets
        };
        if (other.length){
            wordObj.otherChars = other;
        }
        wordObjArr.push(wordObj);
    } catch (e) {
        print(e);
        print(word);
    }
}
//printjson(wordObjArr)
db = connect("localhost/words");
db.word_stats.drop();
db.word_stats.ensureIndex({word: 1}, {unique: true});
db.word_stats.insert(wordObjArr);