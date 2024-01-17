import { API_LIST } from "~/api/config";

export function useCustomRouter() {
  const router = useRouter();
  const route = useRoute();

  function getPageNo(page_no) {
    let page_no_temp = 1;

    if (page_no) {
      page_no_temp = parseInt(page_no, 10);
    } else {
      page_no_temp = parseInt(route.query.page_no ?? 1, 10);
    }

    console.log("page_no_temp", page_no_temp);

    const pageNo = typeof page_no_temp === "number" && page_no_temp >= 1 ? page_no_temp : 1;
    console.log("pageNo", pageNo);
    return pageNo;
  }

  function goPage(page_no = 1, q = "") {
    if (q === "") {
      q = route.query.q ?? "";
    }

    router.push({
      path: "/search",
      query: {
        q,
        page_no: getPageNo(page_no),
      },
    });
  }

  function routeChangeListener(fn, { immediate = true, deep = true } = {}) {
    watch(
      () => route.fullPath,
      (to) => {
        console.log("useCustomRouter routeChangeListener to", to);

        if (typeof fn === "function") {
          fn(to);
        }
      },
      {
        immediate,
        deep,
      },
    );
  }

  return {
    goPage,
    getPageNo,
    routeChangeListener,
  };
}

export async function useFetchComplete(q = "") {
  const { data, pending, error, refresh } = await useFetch(`${API_LIST.complete}`, {
    query: {
      q,
    },
  });

  return data;
}

export async function useFetchSearch(q = "", page_no = 1) {
  const { data, pending, error, refresh } = await useFetch(`${API_LIST.search}`, {
    query: {
      q,
      page_no,
    },
  });

  return data;
}
