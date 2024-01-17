import { findAll } from "highlight-words-core";

/**
 * Highlighted search words
 * @param {string} textToHighlight
 * @param {string} searchWords
 * @returns
 * @example
 * const textToHighlight = "This is some text to highlight.";
 * const searchWords = ["This", "i"];
 */
export function getHighlightedText(textToHighlight, searchWords) {
  // 'dsaf dsaf    dsafasfas  dfsafasf'.split(/\s+/g); ==> ['dsaf', 'dsaf', 'dsafasfas', 'dfsafasf']
  // ''.split(/\s+/g); ==> ['']
  const searchWordsArr = searchWords.split(/\s+/g);

  const chunks = findAll({
    searchWords: searchWordsArr,
    textToHighlight,
    autoEscape: true,
  });

  const highlightedText = chunks
    .map((chunk) => {
      const { end, highlight, start } = chunk;
      const text = textToHighlight.substr(start, end - start);
      if (highlight) {
        return `<em>${text}</em>`;
      } else {
        return text;
      }
    })
    .join("");

  return highlightedText;
}
