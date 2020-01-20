<template>
  <div ref="ide" class="h-100">
    <v-container fluid class="h-100 pa-0">
      <v-row class="h-100">
        <v-col cols="2" class="h-100 pa-0 pl-3">
          <FileTree
            :items="filetree"
            :load-children="fetchChildren"
            :style="filetreeStyle"
            class="file-tree"
            @click-item="clickFileTreeItem"
          ></FileTree>
        </v-col>
        <v-col cols="10" class="h-100 pa-0">
          <div id="editor" class="h-100"></div>
          <div class="console-view">
            <div class="console-toolbar">
              <v-row>
                <v-col cols="3">
                  <v-menu offset-y>
                    <template v-slot:activator="{on}">
                      <v-btn text v-on="on">
                        <v-icon>mdi-console</v-icon>
                        <v-icon>mdi-menu-down</v-icon>
                        <span>{{consoleIndex}}</span>
                      </v-btn>
                    </template>
                    <v-list>
                      <v-list-item
                        v-for="index in consoleCount"
                        :key="index"
                        @click="consoleIndex = index"
                      >
                        <v-list-item-title>{{index}}</v-list-item-title>
                      </v-list-item>
                    </v-list>
                  </v-menu>
                </v-col>
                <v-col cols="1" offset="8">
                  <v-btn icon @click="appendConsole">
                    <v-icon>mdi-plus</v-icon>
                  </v-btn>
                </v-col>
              </v-row>
            </div>
            <div class="console-container">
              <iframe
                v-for="index in consoleCount"
                :key="index"
                :src="consoleURL"
                :style="consoleStyle(index)"
                class="w-100 h-100 console"
              ></iframe>
            </div>
          </div>
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>

<script>
import _ from "lodash";
import ajax from "@/assets/js/ajax.js";
import FileTree from "@/components/FileTree.vue";
import { Url, urlFiles, urlFolders, urlLessons } from "@/assets/js/url.js";

const extension = str => {
  const components = str.split(".");
  const extension =
    components.length == 0 ? "" : components[components.length - 1];
  return extension;
};

const aceMode = path => {
  const ext = extension(path);
  const modes = {
    js: "javascript",
    php: "php"
  };
  const base = "ace/mode/";
  if (!modes[ext]) {
    return base + "text";
  }
  return base + modes[ext];
};

const fileIcon = name => {
  const ext = extension(name);
  if (
    ext == "bash_history" ||
    ext == "bash_logout" ||
    ext == "bash_profile" ||
    ext == "bashrc"
  ) {
    return "mdi-bash";
  } else if (ext == "php") {
    return "mdi-language-php";
  }
  return "mdi-note-text-outline";
};

const makeChildrenFromResponse = response => {
  const children = [];
  if (!response.data.ChildFolders) {
    response.data.ChildFolders = [];
  }
  if (!response.data.ChildFiles) {
    response.data.ChildFiles = [];
  }
  return response.data.ChildFolders.concat(response.data.ChildFiles).map(
    child => {
      const c = {
        name: child.Name,
        path: child.Path
      };
      if (child.hasOwnProperty("Text")) {
        c.file = true;
        c.icon = fileIcon(child.Name);
      } else {
        c.file = false;
        c.open = false;
        c.icon = "mdi-folder";
        c.children = [];
      }
      return c;
    }
  );
};

let editor = null;
export default {
  layout: "lesson",
  components: {
    FileTree
  },
  data() {
    return {
      consoleCount: 1,
      consoleIndex: 1,
      delayedUpdate: _.debounce(this.updateFileText, 1000),
      file: null,
      filetree: [],
      lesson: null,
      ideHeight: 0
    };
  },
  computed: {
    consoleURL() {
      return this.lesson
        ? process.env.serverHost + ":" + this.lesson.HostConsolePort
        : "";
    },
    filetreeStyle() {
      return {
        height: this.ideHeight * 0.99 + "px"
      };
    }
  },
  created() {
    const lessonID = this.$route.params.id;
    const urlL = new Url(urlLessons);
    const lessonData = {
      id: lessonID
    };
    ajax.get(urlL.base, lessonData).then(response => {
      this.lesson = response.data[0];
    });
    const urlF = new Url(urlFolders);
    const folderData = {
      lessonID,
      path: "/"
    };
    ajax.get(urlF.base, folderData).then(response => {
      this.filetree = makeChildrenFromResponse(response);
    });
  },
  mounted() {
    this.ideHeight = this.$refs.ide.clientHeight;
    this.setupAce();
  },
  methods: {
    appendConsole() {
      ++this.consoleCount;
      this.consoleIndex = this.consoleCount;
    },
    clickFileTreeItem(item) {
      if (item.file) {
        const url = new Url(urlFiles);
        const data = {
          lessonID: this.lesson.ID,
          path: item.path
        };
        ajax.get(url.base, data).then(response => {
          item.text = response.data.Text;
          this.file = item;
          editor.setReadOnly(false);
          editor.setValue(item.text);
          editor.getSession().setMode(aceMode(item.path));
        });
      } else {
        item.open = !item.open;
        item.icon = item.open ? "mdi-folder-open" : "mdi-folder";
        if (!item.open) {
          item.children = [];
        }
      }
    },
    consoleStyle(index) {
      return index === this.consoleIndex ? { "z-index": 1 } : {};
    },
    fetchChildren(item) {
      if (item.file) {
        return;
      }
      const url = new Url(urlFolders);
      const data = {
        lessonID: this.$route.params.id,
        path: item.path
      };
      ajax.get(url.base, data).then(response => {
        item.children = makeChildrenFromResponse(response);
      });
    },
    setupAce() {
      ace.edit("editor");
      editor = ace.edit("editor");
      editor.$blockScrolling = Infinity;
      editor.setOptions({
        enableBasicAutocompletion: true,
        enableSnippets: false,
        enableLiveAutocompletion: true
      });
      editor.setTheme("ace/theme/xcode");
      editor.setFontSize(20);
      editor.setReadOnly(true);
      editor.session.on("change", delta => {
        if (delta.action === "insert" && this.file.text !== editor.getValue()) {
          this.file.text = editor.getValue();
          this.delayedUpdate();
        }
      });
    },
    updateFileText() {
      const url = new Url(urlFiles);
      const data = {
        lessonID: this.lesson.ID,
        path: this.file.path,
        text: this.file.text
      };
      ajax.put(url.base, data);
    }
  }
};
</script>

<style>
#editor {
  height: 60%;
}
.console {
  position: absolute;
  left: 0;
  top: 0;
}
.console-container {
  height: 80%;
  position: relative;
}
.console-toolbar {
  height: 20%;
}
.console-view {
  border-top: 1px solid lightgrey;
  height: 39%;
  position: relative;
}
.file-tree {
  overflow: scroll;
  white-space: nowrap;
}
</style>
