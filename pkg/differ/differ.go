package differ

// CalculateDelta returns list of changes found between source and destination chunk maps
func CalculateDelta(src, dst ChunkMap, includeUnchangedChunks bool) Delta {
	changes := make(Delta, 0)

	calculateSourceChanges(&src, &dst, &changes, includeUnchangedChunks)
	calculateDestinationChanges(&src, &dst, &changes)

	return changes
}

func calculateSourceChanges(src, dst *ChunkMap, changes *Delta, includeUnchangedChunks bool) *Delta {
	for srcSum, srcChunk := range *src {
		dstChunk, exists := (*dst)[srcSum]

		if exists {
			indicesEqual := srcChunk.Index == dstChunk.Index

			if indicesEqual && includeUnchangedChunks {
				change := Change{
					Operation: Unchanged,
					From:      srcChunk.Index,
					Bytes:     srcChunk.Bytes,
				}
				*changes = append(*changes, change)
			}

			if !indicesEqual {
				change := Change{
					Operation: Move,
					From:      srcChunk.Index,
					To:        dstChunk.Index,
					Bytes:     srcChunk.Bytes,
				}
				*changes = append(*changes, change)
			}

			continue
		}

		change := Change{
			Operation: Delete,
			From:      srcChunk.Index,
			Bytes:     srcChunk.Bytes,
		}
		*changes = append(*changes, change)
	}

	return changes
}

func calculateDestinationChanges(src, dst *ChunkMap, changes *Delta) *Delta {
	for dstSum, dstChunk := range *dst {
		_, exists := (*src)[dstSum]
		if exists {
			continue
		}

		change := Change{
			Operation: Add,
			To:        dstChunk.Index,
			Bytes:     dstChunk.Bytes,
		}
		*changes = append(*changes, change)
	}

	return changes
}
