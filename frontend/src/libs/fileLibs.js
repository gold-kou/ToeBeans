export function isCorrectSize(file) {
  /** @type {number} 許容する最大サイズ(10MB). */
  var maxSize = 1024 * 1024 * 10;
  return file.size <= maxSize;
}

export function getBase64(file, cb) {
  let reader = new FileReader();
  reader.readAsDataURL(file);
  reader.onload = function () {
    const r = reader.result;
    cb(r.slice(r.indexOf(",") + 1));
  };
  reader.onerror = function (error) {
    throw error;
  };
}
