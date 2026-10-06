/* eslint-disable no-restricted-globals */
/* eslint-disable no-underscore-dangle */
/* eslint-disable no-console */

import { precacheAndRoute } from 'workbox-precaching';

precacheAndRoute(self.__WB_MANIFEST);

const dataCacheName = 'trudgeData-v1.0.0001';
const cacheName = 'trudgeCache';

let filesToCache = [
  '/',
];

self.addEventListener('install', (e) => {
  self.skipWaiting();

  if (typeof self.__precacheManifest !== 'undefined') {
    const manifestCache = [];
    self.__precacheManifest.forEach((v) => {
      if (v.url !== '/index.html') {
        manifestCache.push(v.url);
      }
    });
    filesToCache = filesToCache.concat(manifestCache);
  } else {
    console.log('[ServiceWorker] Failed to load Manifest');
  }

  e.waitUntil(
    caches.open(cacheName).then((cache) => {
      fetch('/').then((response) => {
        cache.put('sw-offline-content', response.clone());
      });
      return cache.addAll(filesToCache);
    }).catch((err) => {
      console.log('[ServiceWorker] Install Fetch Error');
      console.log(err);
    }),
  );
});

self.addEventListener('activate', (e) => {
  e.waitUntil(
    caches.keys().then((keyList) => Promise.all(
      keyList.map((key) => (
        (key !== cacheName && key !== dataCacheName) ? caches.delete(key) : null
      )),
    )),
  );
  /*
   * Fixes a corner case in which the app wasn't returning the latest data.
   * You can reproduce the corner case by commenting out the line below and
   * then doing the following steps: 1) load app for first time so that the
   * initial New York City data is shown 2) press the refresh button on the
   * app 3) go offline 4) reload the app. You expect to see the newer NYC
   * data, but you actually see the initial data. This happens because the
   * service worker is not yet activated. The code below essentially lets
   * you activate the service worker faster.
   */
  return self.clients.claim();
});

self.addEventListener('fetch', (e) => {
  const dataUrl = 'trudge.furzr.com/api/';
  const logoffUrl = 'trudge.furzr.com/auth/logout';
  if (e.request.url.indexOf(dataUrl) > -1 || e.request.url.indexOf(logoffUrl) > -1) {
    /*
     * When the request URL contains dataUrl, the app is asking for fresh
     * weather data. In this case, the service worker always goes to the
     * network and then caches the response. This is called the "Cache then
     * network" strategy:
     * https://jakearchibald.com/2014/offline-cookbook/#cache-then-network
     */
    e.respondWith(
      caches.open(dataCacheName).then((cache) => fetch(e.request).then((response) => {
        cache.put(e.request.url, response.clone()).catch(() => {
          // Don't care
        });
        return response;
      })).catch((err) => {
        console.log('[ServiceWorker] Fetch - API Failure');
        console.log(err);
      }),
    );
  } else {
    /*
     * The app is asking for app shell files. In this scenario the app uses the
     * "Cache, falling back to the network" offline strategy:
     * https://jakearchibald.com/2014/offline-cookbook/#cache-falling-back-to-network
     */
    e.respondWith(
      caches.match(e.request).then((response) => response || fetch(e.request).catch(() => ((e.request.mode === 'navigate'
        || (e.request.method === 'GET' && e.request.headers.get('accept').includes('text/html')))
        ? caches.match('sw-offline-content').then((offline) => offline) : null)))
        .catch((err) => {
          console.log('[ServiceWorker] Fetch - Static Failure');
          console.log(err);
        }),
    );
  }
});

self.addEventListener('message', (e) => {
  console.log('SW Message Received');
  console.log(e);
});

self.addEventListener('push', (e) => {
  console.log('Push Notice');
  console.log(e);
  console.log(document);
});

self.addEventListener('notificationclick', (event) => {
  event.notification.close();

  // eslint-disable-next-line no-undef
  event.waitUntil(clients.matchAll({
    type: 'window',
  }).then((clientList) => {
    clientList.forEach((client) => {
      if ('focus' in client) {
        console.log('postMessage');
        client.postMessage({ event: event.action });
        client.focus();
      }
    });
  }));
});
