export const IS_PROD = import.meta.env.PROD;
export const ENV = import.meta.env.MODE;

// true 'production' or false 'development'
console.log(IS_PROD, ENV);

// TODO!!! production host
// default: domain/api/
// @example:
// when your search site is https://my-search.com, then your api is https://my-search.com/api/search and https://my-search.com/api/complete
export const HOST = ENV === "development" ? "http://localhost:9999/api" : "/api";

export const API_LIST = {
  search: `${HOST}/search`,
  complete: `${HOST}/complete`,
};
