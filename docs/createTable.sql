CREATE TABLE [dbo].[OpsTerminalInfo](
	[RegionId] [nvarchar](6) NOT NULL,
	[StoreId] [nvarchar](6) NOT NULL,
	[TerminalId] [nvarchar](20) NOT NULL,
	[LastActionDateTime] [datetime] NULL
) ON [PRIMARY]