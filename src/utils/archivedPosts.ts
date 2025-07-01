import { getAllSortedPosts } from "./posts";

const sortedPosts = await getAllSortedPosts();

export const archivedPostsByYear = sortedPosts.reduce((acc, post) => {
  const year = post.data.date.getFullYear();
  if (!acc[year]) {
    acc[year] = [];
  }
  acc[year].push(post);
  return acc;
}, {} as Record<number, typeof sortedPosts>);

export const archivedYears = Object.keys(archivedPostsByYear).sort((a, b) => Number(b) - Number(a));