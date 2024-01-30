<template>
  <form id="search" method="GET" action="/search" role="search">
    <div class="search-view">
      <a v-if="props.isLogoShow" class="search-logo" href="/" tabindex="0">
        <svg width="92mm" height="92mm" viewBox="0 0 92 92">
          <g transform="translate(-40.921 -17.417)">
            <circle
              cx="75.921"
              cy="53.903"
              r="30"
              fill="none"
              fill-opacity="1"
              stroke="#3050ff"
              stroke-width="10"
              stroke-miterlimit="4"
              stroke-dasharray="none"
              stroke-opacity="1"
            ></circle>
            <path
              d="M67.515 37.915a18 18 0 0 1 21.051 3.313 18 18 0 0 1 3.138 21.078"
              fill="none"
              fill-opacity="1"
              stroke="#3050ff"
              stroke-width="5"
              stroke-miterlimit="4"
              stroke-dasharray="none"
              stroke-opacity="1"
            ></path>
            <rect
              width="18.846"
              height="39.963"
              x="3.706"
              y="122.09"
              ry="0"
              transform="rotate(-46.235)"
              opacity="1"
              fill="#3050ff"
              fill-opacity="1"
              stroke="none"
              stroke-width="8"
              stroke-miterlimit="4"
              stroke-dasharray="none"
              stroke-opacity="1"
            ></rect>
          </g>
        </svg>
      </a>

      <div class="search-box">
        <input
          ref="search_input"
          @focus="searchInputFocusHandler"
          @blur="searchInputBlurHandler"
          v-model="searchValue"
          id="q"
          name="q"
          type="text"
          placeholder="search..."
          autocomplete="off"
          autocapitalize="none"
          spellcheck="false"
          autocorrect="off"
          dir="auto"
        />

        <button
          type="button"
          @click="clearBtnHandler"
          aria-label="clear"
          class="clear-search"
          :class="{ empty: searchValue === '' }"
        >
          <span>
            <svg class="icon-big" viewBox="0 0 512 512" aria-hidden="true">
              <path
                fill="none"
                stroke="currentColor"
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="32"
                d="M368 368L144 144M368 144L144 368"
              ></path>
            </svg>
          </span>
          <span class="hidden">clear</span>
        </button>

        <button class="send-search" type="submit" aria-label="search">
          <span>
            <svg class="icon-big" viewBox="0 0 512 512" aria-hidden="true">
              <path
                d="M221.09 64a157.09 157.09 0 10157.09 157.09A157.1 157.1 0 00221.09 64z"
                fill="none"
                stroke="currentColor"
                stroke-miterlimit="10"
                stroke-width="32"
              ></path>
              <path
                fill="none"
                stroke="currentColor"
                stroke-linecap="round"
                stroke-miterlimit="10"
                stroke-width="32"
                d="M338.29 338.29L448 448"
              ></path>
            </svg>
          </span>
          <span class="hidden">search</span>
        </button>

        <div class="autocomplete" :class="{ open: autocompleteIsShow }">
          <ul>
            <li
              @click="onAutocompleteItemClick(q)"
              :data-autocomplete-value="q"
              v-for="q in autocompleteArr"
            >
              <span v-html="getHighlightedText(q, searchValue)"></span>
            </li>

            <!-- <li data-autocomplete-value="123demands"><strong>123d</strong>emands</li>
            <li data-autocomplete-value="123dough fine foods &amp;amp; provisions">
              <strong>123d</strong>ough fine foods &amp; provisions
            </li>
            <li data-autocomplete-value="123dough"><strong>123d</strong>ough</li>
            <li data-autocomplete-value="123demands pet sim 99"><strong>123d</strong>emands pet sim 99</li>
            <li data-autocomplete-value="123dj"><strong>123d</strong>j</li>
            <li data-autocomplete-value="123d design"><strong>123d</strong> design</li>
            <li data-autocomplete-value="123dentist"><strong>123d</strong>entist</li>
            <li data-autocomplete-value="123d catch"><strong>123d</strong> catch</li>
            <li data-autocomplete-value="123deals"><strong>123d</strong>eals</li>
            <li data-autocomplete-value="123dough bakery"><strong>123d</strong>ough bakery</li> -->
          </ul>
        </div>
      </div>
    </div>
  </form>
</template>

<script setup>
import { debounce } from "lodash-es";
import { useCustomRouter } from "~/utils/hooks";
import { getHighlightedText } from "~/utils";
import { useSearchStore } from "~/stores/search";

const { getSearchResults } = useSearchStore();

const { goPage, routeChangeListener } = useCustomRouter();
const completeCache = {};

const props = defineProps({
  isLogoShow: {
    type: Boolean,
    default: true,
  },
});

const search_input = ref(null);
const searchValue = ref("");
const autocompleteIsShow = ref(false);
const autocompleteArr = ref([]);
const searchValueIsFocus = ref(false);

routeChangeListener(() => {
  searchBoxInit();
});

watch(searchValue, () => {
  if (searchValueIsFocus.value === true) {
    searchDebounceForAutocompleteHandler();
  }
});

function onAutocompleteItemClick(q) {
  goPage(1, q).then((_) => {
    getSearchResults();
  });
}

function searchInputFocusHandler() {
  searchValueIsFocus.value = true;
  autocompleteIsShow.value = true;

  document.body.classList.add("search-input-focus");
}

function searchInputBlurHandler() {
  searchValueIsFocus.value = false;

  // because need click, so delay
  setTimeout(() => {
    autocompleteIsShow.value = false;
    document.body.classList.remove("search-input-focus");
  }, 100);
}

const searchDebounceForAutocompleteHandler = debounce(async () => {
  const q = searchValue.value;

  console.log(q);

  // check autocomplete
  if (q === "") {
    autocompleteIsShow.value = false;
    return;
  }

  let val = [];
  let completeDataQuery = q;

  if (Array.isArray(completeCache[q]) && completeCache[q].length > 0) {
    console.log(q, completeCache[q]);

    val = completeCache[q];
  } else {
    const data = await useFetchComplete(q);

    try {
      console.log("data", data.value);

      completeDataQuery = data.value.query;
      val = data.value.results.slice(0, 10).map((item) => {
        return item.text;
      });
      completeCache[q] = val;
    } catch (err) {
      console.error(err);
    }
  }

  // if ajax delay, so check current search value is ajax data
  if (completeDataQuery === searchValue.value) {
    autocompleteArr.value = val;
  } else {
    autocompleteArr.value = [];
  }

  if (searchValue.value === "") {
    autocompleteIsShow.value = false;
  } else {
    autocompleteIsShow.value = true;
  }
}, 200);

function clearBtnHandler() {
  searchValue.value = "";

  if (search_input.value) {
    search_input.value.focus();
  }
}

function searchBoxInit() {
  const route = useRoute();

  if (route.query && route.query.q) {
    searchValue.value = route.query.q;
  }
}
</script>

<style lang="less" scoped>
.hidden {
  display: none;
}

.search-view {
  display: flex;
  justify-content: center;
  align-items: center;
}

.search-logo {
  margin-right: 0.5em;

  svg {
    flex: 1;
    width: 30px;
    height: 30px;
    margin: 0.5rem 0 auto 0;
  }
}

.search-box {
  border: 1px solid var(--color-search-border);
  box-shadow: var(--color-search-shadow);
  border-radius: 2em;
  width: 46rem;
  display: inline-flex;
  flex-direction: row;
  white-space: nowrap;
  position: relative;
}

.clear-search {
  display: block;
  border-collapse: separate;
  box-sizing: border-box;
  width: 1.8rem;
  margin: 0;
  padding: 1rem 0.2rem;
  background: none repeat scroll 0 0 var(--color-search-background);
  border: none;
  outline: 0;
  color: var(--color-search-font);
  font-size: 1.1rem;

  &:hover {
    color: var(--color-search-background-hover);
  }

  &.empty * {
    display: none;
  }
}

.send-search {
  display: block;
  margin: 0;
  padding: 1rem;
  background: none repeat scroll 0 0 var(--color-search-background);
  border: none;
  outline: 0;
  color: var(--color-search-font);
  font-size: 1.1rem;
  z-index: 2;
  border-radius: 0 2em 2em 0;

  &:hover {
    cursor: pointer;
    background-color: var(--color-search-background-hover);
    color: var(--color-search-background);
  }
}

.icon-big {
  width: 1.5rem;
  height: 1.5rem;
  display: inline-block;
  vertical-align: bottom;
  line-height: 1;
  text-decoration: inherit;
  transform: scale(1, 1);
}

#q {
  width: 100%;
  padding-left: 1.4rem;
  padding-right: 0 !important;
  border-radius: 2em 0 0 2em;
  display: block;
  margin: 0;
  padding: 1rem;
  background: none repeat scroll 0 0 var(--color-search-background);
  border: none;
  outline: 0;
  color: var(--color-search-font);
  font-size: 1.1rem;
  z-index: 2;
}

.autocomplete {
  position: absolute;
  width: 100%;
  left: 0;
  max-height: 0;
  overflow-y: hidden;
  .ltr-text-align-left();

  .rounded-corners;

  &:active,
  &:focus,
  &:hover {
    background-color: var(--color-autocomplete-background);
  }

  &:empty {
    display: none;
  }

  > ul {
    list-style-type: none;
    margin: 0;
    padding: 0;

    > li {
      cursor: pointer;
      padding: 0.5rem 1rem;

      :deep(em) {
        font-weight: bold;
      }

      &.active,
      &:active,
      &:focus,
      &:hover {
        background-color: var(--color-autocomplete-background-hover);

        a:active,
        a:focus,
        a:hover {
          text-decoration: none;
        }
      }

      &.locked {
        cursor: inherit;
      }
    }
  }

  &.open {
    display: block;
    background-color: var(--color-autocomplete-background);
    color: var(--color-autocomplete-font);
    max-height: 20rem;
    overflow-y: auto;
    z-index: 100;
    margin-top: 3.5rem;
    border-radius: 0.8rem;
    box-shadow: 0 2px 8px rgb(34 38 46 / 25%);

    &:empty {
      display: none;
    }
  }
}

@media screen and (max-width: @phone) {
  // .autocomplete {
  //   width: 100%;

  //   >ul>li {
  //     padding: 1rem;
  //   }
  // }
}
</style>
