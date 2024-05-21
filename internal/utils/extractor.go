package utils

func ExtractIDs[T any](items *[]T, getID func(T) string) []string {
    var ids []string
    if items != nil {
        for _, item := range *items {
            ids = append(ids, getID(item))
        }
    }
    return ids
}
