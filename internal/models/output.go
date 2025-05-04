package models

// NestedDigestResults represents the nested structure for container digests
// Format:
// {
//   "docker.io": {
//     "library/alpine": {
//       "3.14": {
//         "linux/amd64": "sha256:71859b0c62df47efaeae4f93698b56a8dddafbf041778fd668bbd1ab45a864f8"
//       }
//     }
//   }
// }
type NestedDigestResults map[string]RepositoryMap

// RepositoryMap maps repository names to their tags
type RepositoryMap map[string]TagMap

// TagMap maps tags to their architectures
type TagMap map[string]ArchMap

// ArchMap maps architectures to their digests
type ArchMap map[string]string
