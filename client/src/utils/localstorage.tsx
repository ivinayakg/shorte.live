const getFromLocalStorage = (key: string) => {
  return JSON.parse(localStorage.getItem(key) ?? "");
};
const setInLocalStorage = (key: string, data: any) => {
  localStorage.setItem(key, JSON.stringify(data));
  return;
};

export { setInLocalStorage, getFromLocalStorage };
