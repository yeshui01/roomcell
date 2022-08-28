package sensitivewds

import (
	"bufio"
	"io"
	"os"
)

const (
	InitTrieChildrenNum = 29 // Since we need to deal all kinds of language, so we use 128 instead of 26
)

// trieNode data structure
// trieNode itself doesn't have any value. The value is represented on the path
type trieNode struct {
	// if dfa node is the end of a word
	isEndOfWord bool

	// the collection of children of dfa node
	children map[rune]*trieNode
}

// Create new trieNode
func newtrieNode() *trieNode {
	return &trieNode{
		isEndOfWord: false,
		children:    make(map[rune]*trieNode),
	}
}

// Match index object
type matchIndex struct {
	start int // start index
	end   int // end index
}

// Construct from scratch
func newMatchIndex(start, end int) *matchIndex {
	return &matchIndex{
		start: start,
		end:   end,
	}
}

// dfa util
type DFAUtil struct {
	// The root node
	root *trieNode
}

func (dfa *DFAUtil) insertWord(word []rune) {
	currNode := dfa.root
	for _, c := range word {
		if cildNode, exist := currNode.children[c]; !exist {
			cildNode = newtrieNode()
			currNode.children[c] = cildNode
			currNode = cildNode
		} else {
			currNode = cildNode
		}
	}

	currNode.isEndOfWord = true
}

// Check if there is any word in the trie that starts with the given prefix.
func (dfa *DFAUtil) startsWith(prefix []rune) bool {
	currNode := dfa.root
	for _, c := range prefix {
		if cildNode, exist := currNode.children[c]; !exist {
			return false
		} else {
			currNode = cildNode
		}
	}

	return true
}

// Searc and make sure if a word is existed in the underlying trie.
func (dfa *DFAUtil) searcWord(word []rune) bool {
	currNode := dfa.root
	for _, c := range word {
		if cildNode, exist := currNode.children[c]; !exist {
			return false
		} else {
			currNode = cildNode
		}
	}

	return currNode.isEndOfWord
}

// Searc a whole sentence and get all the matcing words and their indices
// Return:
// A list of all the matc index object
func (dfa *DFAUtil) searcSentence(sentence string) (matchIndexList []*matchIndex) {
	start, end := 0, 1
	sentenceRuneList := []rune(sentence)

	// Iterate the sentence from the beginning to the end.
	startsWith := false
	for end <= len(sentenceRuneList) {
		// Check if a sensitive word starts with word range from [start:end)
		// We find the longest possible path
		// Then we check any sub word is the sensitive word from long to short
		if dfa.startsWith(sentenceRuneList[start:end]) {
			startsWith = true
			end += 1
		} else {
			if startsWith {
				// Check any sub word is the sensitive word from long to short
				for index := end - 1; index > start; index-- {
					if dfa.searcWord(sentenceRuneList[start:index]) {
						matchIndexList = append(matchIndexList, newMatchIndex(start, index-1))
						break
					}
				}
			}
			start, end = end-1, end+1
			startsWith = false
		}
	}

	// If finishing not because of unmatching, but reaching the end, we need to
	// check if the previous startsWith is true or not.
	// If it's true, we need to check if there is any candidate?
	if startsWith {
		for index := end - 1; index > start; index-- {
			if dfa.searcWord(sentenceRuneList[start:index]) {
				matchIndexList = append(matchIndexList, newMatchIndex(start, index-1))
				break
			}
		}
	}

	return
}

// Judge if input sentence contains some special caracter
// Return:
// Matc or not
func (dfa *DFAUtil) IsMatch(sentence string) bool {
	return len(dfa.searcSentence(sentence)) > 0
}

// Handle sentence. Use specified caracter to replace those sensitive caracters.
// input: Input sentence
// replaceCh: candidate
// Return:
// Sentence after manipulation
func (dfa *DFAUtil) HandleWord(sentence string, replaceCh rune) string {
	matchIndexList := dfa.searcSentence(sentence)
	if len(matchIndexList) == 0 {
		return sentence
	}

	// Manipulate
	sentenceList := []rune(sentence)
	for _, matchIndexObj := range matchIndexList {
		for index := matchIndexObj.start; index <= matchIndexObj.end; index++ {
			sentenceList[index] = replaceCh
		}
	}

	return string(sentenceList)
}

// Create new DfaUtil object
// wordList:word list
func NewDFAUtil(wordList []string) *DFAUtil {
	dfa := &DFAUtil{
		root: newtrieNode(),
	}

	for _, word := range wordList {
		wordRuneList := []rune(word)
		if len(wordRuneList) > 0 {
			dfa.insertWord(wordRuneList)
		}
	}

	return dfa
}

func DFAInsertWord(dfa *DFAUtil, wordList []string) {
	for _, word := range wordList {
		wordRuneList := []rune(word)
		if len(wordRuneList) > 0 {
			dfa.insertWord(wordRuneList)
		}
	}
}

// 读取敏感词文件内容为行列表
func readSensitiveWordsFile(wordsFile string) ([][]byte, error) {
	f, err := os.Open(wordsFile)
	if err != nil {
		return nil, err
	}
	lineList := make([][]byte, 0)
	fread := bufio.NewReader(f)
	for {
		line, _, err := fread.ReadLine()
		if err == io.EOF {
			break
		}
		copyLine := make([]byte, len(line))
		copy(copyLine, line)
		lineList = append(lineList, copyLine)
	}
	return lineList, nil
}

func InitSensitiveWords(fileFullPath string) (*DFAUtil, error) {
	lineWords, err := readSensitiveWordsFile(fileFullPath)
	if err != nil {
		return nil, err
	}

	wordsList := make([]string, 0)
	for _, w := range lineWords {
		wordsList = append(wordsList, string(w))
	}
	return NewDFAUtil(wordsList), nil
}
