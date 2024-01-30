import { defineStore } from "pinia";
import { useFetchSearch } from "~/utils/hooks";
const { getPageNo } = useCustomRouter();

export const useSearchStore = defineStore("searchStore", () => {
  const searchStore = reactive({
    results: [],
    loading: false,
    total: 100,
    q: "",
  });

  async function getSearchResults() {
    console.log("call getSearchResults", searchStore);

    if (searchStore.loading) {
      return;
    }
    searchStore.loading = true;

    const route = useRoute();

    searchStore.results.length = 0;
    searchStore.q = route.query.q ?? "";

    const page_no = getPageNo();
    const data = await useFetchSearch(searchStore.q, page_no);

    try {
      console.log("data", data, data.value, data.value.results, typeof data.value.results);

      if (data?.value?.results) {
        searchStore.results = data.value.results;

        // no total, so set 100
        searchStore.total = 100;
      }
    } catch (err) {
      console.error(err);
    }

    searchStore.loading = false;
  }

  return {
    searchStore,
    getSearchResults,
  };
});
