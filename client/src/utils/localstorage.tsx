const getFromLocalStorage = (key: string) => {
  const storedValue = localStorage.getItem(key);
  // Check if the stored value is null or undefined
  if (storedValue === null || storedValue === undefined) {
    return null; // Or any default value you want to return
  }

  try {
    return JSON.parse(storedValue);
  } catch (error) {
    console.error("Error parsing JSON from localStorage:", error);
    return null; // Or handle the error in a way that makes sense for your application
  }
};

const setInLocalStorage = (key: string, data: any) => {
  localStorage.setItem(key, JSON.stringify(data));
  return;
};

export { setInLocalStorage, getFromLocalStorage };
