module.exports = {
    theme: 'cosmos',
    title: 'Merlion Documentation',
    locales: {
        '/': {
            lang: 'en-US'
        },
    },
    markdown: {
        extendMarkdown: (md) => {
            md.use(require("markdown-it-katex"));
        },
    },
    head: [
        [
            "link",
            {
                rel: "stylesheet",
                href:
                    "https://cdnjs.cloudflare.com/ajax/libs/KaTeX/0.5.1/katex.min.css",
            },
        ],
        [
            "link",
            {
                rel: "stylesheet",
                href:
                    "https://cdn.jsdelivr.net/github-markdown-css/2.2.1/github-markdown.css",
            },
        ],
    ],
    plugins: [
        'vuepress-plugin-element-tabs'
    ],
    themeConfig: {
        repo: 'merlion-zone/merlion',
        docsRepo: 'merlion-zone/merlion',
        docsBranch: 'main',
        docsDir: 'docs',
        editLinks: true,
        custom: true,
    }
}