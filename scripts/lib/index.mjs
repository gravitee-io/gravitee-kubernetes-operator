const { log } = console;

const newLoggerFn =
  (color) =>
  (...args) =>
    log(color(...args));

const $Quote = $.quote;
const $NoQuote = (unescaped) => unescaped;

export async function time(fn) {
  const start = Date.now();
  await fn();
  const end = Date.now();
  green(`Done in ${(end - start) / 1000}s`);
}

export function toggleVerbosity(verbose = false) {
  $.verbose = verbose;
}

export function setNoQuoteEscape() {
    $.quote = $NoQuote;
}

export function setQuoteEscape() {
    $.quote = $Quote;
}

export const blue = newLoggerFn(chalk.blue);
export const yellow = newLoggerFn(chalk.yellow);
export const magenta = newLoggerFn(chalk.magenta);
export const green = newLoggerFn(chalk.green);
