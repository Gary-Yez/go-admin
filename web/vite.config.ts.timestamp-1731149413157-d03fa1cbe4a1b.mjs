// vite.config.ts
import { defineConfig } from "file:///F:/%E6%88%91%E7%9A%84%E9%A1%B9%E7%9B%AE/%E8%87%AA%E7%94%A8%E5%90%8E%E5%8F%B0%E6%A1%86%E6%9E%B6/admin-web/node_modules/vite/dist/node/index.js";
import vue from "file:///F:/%E6%88%91%E7%9A%84%E9%A1%B9%E7%9B%AE/%E8%87%AA%E7%94%A8%E5%90%8E%E5%8F%B0%E6%A1%86%E6%9E%B6/admin-web/node_modules/@vitejs/plugin-vue/dist/index.mjs";
import AutoImport from "file:///F:/%E6%88%91%E7%9A%84%E9%A1%B9%E7%9B%AE/%E8%87%AA%E7%94%A8%E5%90%8E%E5%8F%B0%E6%A1%86%E6%9E%B6/admin-web/node_modules/unplugin-auto-import/dist/vite.js";
import Components from "file:///F:/%E6%88%91%E7%9A%84%E9%A1%B9%E7%9B%AE/%E8%87%AA%E7%94%A8%E5%90%8E%E5%8F%B0%E6%A1%86%E6%9E%B6/admin-web/node_modules/unplugin-vue-components/dist/vite.js";
import { ElementPlusResolver } from "file:///F:/%E6%88%91%E7%9A%84%E9%A1%B9%E7%9B%AE/%E8%87%AA%E7%94%A8%E5%90%8E%E5%8F%B0%E6%A1%86%E6%9E%B6/admin-web/node_modules/unplugin-vue-components/dist/resolvers.js";
var vite_config_default = defineConfig({
  plugins: [
    vue(),
    AutoImport({
      resolvers: [ElementPlusResolver()]
    }),
    Components({
      resolvers: [ElementPlusResolver()]
    })
  ]
});
export {
  vite_config_default as default
};
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsidml0ZS5jb25maWcudHMiXSwKICAic291cmNlc0NvbnRlbnQiOiBbImNvbnN0IF9fdml0ZV9pbmplY3RlZF9vcmlnaW5hbF9kaXJuYW1lID0gXCJGOlxcXFxcdTYyMTFcdTc2ODRcdTk4NzlcdTc2RUVcXFxcXHU4MUVBXHU3NTI4XHU1NDBFXHU1M0YwXHU2ODQ2XHU2N0I2XFxcXGFkbWluLXdlYlwiO2NvbnN0IF9fdml0ZV9pbmplY3RlZF9vcmlnaW5hbF9maWxlbmFtZSA9IFwiRjpcXFxcXHU2MjExXHU3Njg0XHU5ODc5XHU3NkVFXFxcXFx1ODFFQVx1NzUyOFx1NTQwRVx1NTNGMFx1Njg0Nlx1NjdCNlxcXFxhZG1pbi13ZWJcXFxcdml0ZS5jb25maWcudHNcIjtjb25zdCBfX3ZpdGVfaW5qZWN0ZWRfb3JpZ2luYWxfaW1wb3J0X21ldGFfdXJsID0gXCJmaWxlOi8vL0Y6LyVFNiU4OCU5MSVFNyU5QSU4NCVFOSVBMSVCOSVFNyU5QiVBRS8lRTglODclQUElRTclOTQlQTglRTUlOTAlOEUlRTUlOEYlQjAlRTYlQTElODYlRTYlOUUlQjYvYWRtaW4td2ViL3ZpdGUuY29uZmlnLnRzXCI7aW1wb3J0IHsgZGVmaW5lQ29uZmlnIH0gZnJvbSAndml0ZSdcbmltcG9ydCB2dWUgZnJvbSAnQHZpdGVqcy9wbHVnaW4tdnVlJ1xuaW1wb3J0IEF1dG9JbXBvcnQgZnJvbSAndW5wbHVnaW4tYXV0by1pbXBvcnQvdml0ZSdcbmltcG9ydCBDb21wb25lbnRzIGZyb20gJ3VucGx1Z2luLXZ1ZS1jb21wb25lbnRzL3ZpdGUnO1xuaW1wb3J0IHsgRWxlbWVudFBsdXNSZXNvbHZlciB9IGZyb20gJ3VucGx1Z2luLXZ1ZS1jb21wb25lbnRzL3Jlc29sdmVycyc7XG4vLyBodHRwczovL3ZpdGUuZGV2L2NvbmZpZy9cbmV4cG9ydCBkZWZhdWx0IGRlZmluZUNvbmZpZyh7XG4gIHBsdWdpbnM6IFtcbiAgICB2dWUoKSxcbiAgICBBdXRvSW1wb3J0KHtcbiAgICAgIHJlc29sdmVyczpbRWxlbWVudFBsdXNSZXNvbHZlcigpXVxuICAgIH0pLFxuICAgIENvbXBvbmVudHMoe1xuICAgICAgcmVzb2x2ZXJzOiBbRWxlbWVudFBsdXNSZXNvbHZlcigpXSxcbiAgICB9KSxcbiAgXSxcbn0pXG4iXSwKICAibWFwcGluZ3MiOiAiO0FBQWtWLFNBQVMsb0JBQW9CO0FBQy9XLE9BQU8sU0FBUztBQUNoQixPQUFPLGdCQUFnQjtBQUN2QixPQUFPLGdCQUFnQjtBQUN2QixTQUFTLDJCQUEyQjtBQUVwQyxJQUFPLHNCQUFRLGFBQWE7QUFBQSxFQUMxQixTQUFTO0FBQUEsSUFDUCxJQUFJO0FBQUEsSUFDSixXQUFXO0FBQUEsTUFDVCxXQUFVLENBQUMsb0JBQW9CLENBQUM7QUFBQSxJQUNsQyxDQUFDO0FBQUEsSUFDRCxXQUFXO0FBQUEsTUFDVCxXQUFXLENBQUMsb0JBQW9CLENBQUM7QUFBQSxJQUNuQyxDQUFDO0FBQUEsRUFDSDtBQUNGLENBQUM7IiwKICAibmFtZXMiOiBbXQp9Cg==
