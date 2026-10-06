<script>
import Axios from 'axios';

export default {
  name: 'CachedRequest',
  data() {
    return {
      currentResponseSource: {},
    };
  },
  methods: {
    getRequest(url, callback = () => {}, failcallback = () => {}) {
      this.cachedRequest(url, callback);
      this.liveRequest(url, callback, failcallback);
    },
    cachedRequest(url, callback) {
      if ('caches' in window) {
        /*
         * Check if the service worker has already cached this city's weather
         * data. If the service worker has the data, then display the cached
         * data while the app fetches the latest data.
         */
        caches.match(url).then((response) => {
          if (response) {
            if (response.status >= 400) {
              return;
            }
            response.json().then((r) => {
              if (typeof(this.currentResponseSource[url]) === 'undefined') {
                callback(r);
              }
            });
          }
        });
      }
    },
    liveRequest(url, callback, failcallback = () => {}) {
      Axios.get(url).then((r) => {
        this.currentResponseSource[url] = 1;
        callback(r.data);
      }).catch((e) => {
        failcallback(e);
      });
    },
  },
};
</script>
