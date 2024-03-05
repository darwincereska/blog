// import rss from '@astrojs/rss';
// import { getCollection } from 'astro:content';
// import config from '@config/config.json';
// import sanitizeHtml from 'sanitize-html';
// import MarkdownIt from 'markdown-it';


// const parser = new MarkdownIt();
// export async function GET(context) {
//   const posts = await getCollection('blog');
//   return rss({
//     title: config.site.title,
//     description: config.site.description,
//     site: context.site,

//     items: posts.map((post) => ({
//       ...post.data,
//       content: sanitizeHtml(parser.render(post.body)),
//       link: `https://blog.durrstudios.dev/blog/${post.slug}/`,

//     })),
//   });
// }

import rss from '@astrojs/rss';
import { getCollection } from 'astro:content';
import config from '@config/config.json';
import sanitizeHtml from 'sanitize-html';
import MarkdownIt from 'markdown-it';


const parser = new MarkdownIt();
export async function GET(context) {
  const posts = await getCollection('blog');
  return rss({
    title: config.site.title,
    description: config.site.description,
    site: 'https://blog.durrstudios.dev',
    
    items: posts.map((post) => ({
      ...post.data,
      content: sanitizeHtml(parser.render(post.body)),
      link: `https://blog.durrstudios.dev/blog/${post.slug}/`,

    })),
  });
}

